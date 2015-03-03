package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type procstart struct {
	UUID      string `json:"uuid"`
	LocalTime string `json:"localtime"`
	Command   string `json:"command"`
	Hostname  string `json:"hostname"`
	IPaddress string `json:"ipaddress"`
	Hash      string `json:"hash"`
}

type heartbeat struct {
	UUID        string `json:"uuid"`
	Ping        string `json:"Ping"`
	RunningTime string `json:"runningtime"`
}

type procend struct {
	UUID        string `json:"uuid"`
	Error       bool   `json:"error"`
	ExitMessage string `json:"exitmessage"`
}

type config struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
}

var uuid = genUuid()
var debug = checkDebugStatus()
var conf = parseConfig()

func main() {

	if len(os.Args) == 1 {
		log.Print("nothing to run")
		os.Exit(1)
	}

	sendCom(os.Args[1:], "/start")
	go sendHeartbeat()
	args := []string(os.Args[2:])
	cmd := exec.Command(os.Args[1], args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	cmd.Start()

	cmd.Wait()

	sendEnd()

}

func checkDebugStatus() bool {

	if os.Getenv("DEBUG") == "TRUE" {
		return true
	}

	return false

}

func genUuid() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return ""
	}

	var fixed string
	fixed = strings.Replace(string(out), "\n", "", -1)

	return fixed
}

func runningTime(startingTime time.Time, endingTime time.Time) string {

	var duration time.Duration = endingTime.Sub(startingTime)

	durationMs := int64(duration / 1000000)
	durationMin := durationMs / 1000 / 60 % 60
	durationHr := durationMs / 1000 / 60 / 60 % 24

	if durationMin > 60 {
		return strconv.FormatInt(durationHr, 10) + " hr"
	}

	return strconv.FormatInt(durationMin, 10) + " min"

}

func sendCom(cmd []string, endpoint string) error {

	// Get IP
	ip, _ := getIP()

	// Get Command and Args
	cmdstring := strings.Join(cmd, " ")

	// Get Current Time
	curTime := string(time.Now().Format(time.RFC3339))

	// Get Hostname
	hostname, _ := os.Hostname()

	// Create a hash of the command, ip and hostname strings
	hashOfInfo := md5String(cmdstring + ip + hostname)

	blob := procstart{
		UUID:      uuid,
		LocalTime: curTime,
		Command:   cmdstring,
		Hostname:  hostname,
		IPaddress: ip,
		Hash:      hashOfInfo,
	}

	jsonBlob, err := json.Marshal(blob)
	if err != nil {
		return err
	}

	err = postJSON(endpoint, jsonBlob)
	if err != nil {
		log.Print("Failed to connect to MotherShip")
	}

	return nil
}

func sendHeartbeat() {

	startingTime := time.Now().UTC()

	c := time.Tick(60 * time.Second)
	for _ = range c {

		endingTime := time.Now().UTC()

		runningTime := runningTime(startingTime, endingTime)

		blob := heartbeat{
			Ping:        string(time.Now().Format(time.RFC3339)),
			UUID:        uuid,
			RunningTime: runningTime,
		}

		endpoint := "/heartbeat"
		jsonBlob, err := json.Marshal(blob)
		if err != nil {
			log.Print(err)
		}

		err = postJSON(endpoint, jsonBlob)
		if err != nil {
			log.Print("Failed to connect to MotherShip")
		}

	}

}

func catchEnd(c <-chan os.Signal) {

	for sig := range c {
		log.Print("Got Ctrl+c")
		os.Exit(1)
	}

}

func sendEnd() {

	blob := procend{
		UUID:        uuid,
		Error:       true,
		ExitMessage: "Fail",
	}

	endpoint := "/end"
	jsonBlob, err := json.Marshal(blob)
	if err != nil {
		log.Print(err)
	}

	err = postJSON(endpoint, jsonBlob)
	if err != nil {
		log.Print("Failed to connect to MotherShip")
	}

}

func parseConfig() *config {

	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Print(err)
	}

	var conf config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Print(err)
	}
	return &conf

}

func postJSON(endpoint string, jsonBlob []byte) error {

	url := "http://" + conf.Hostname + ":" + conf.Port + endpoint
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBlob))

	req.Header.Set("X-Custom-Header", "MyServer")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil

}

func getIP() (string, error) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", nil
	}

	var IP string

	for _, add := range addrs {
		if ipadd, ok := add.(*net.IPNet); ok && !ipadd.IP.IsLoopback() {
			// Add get ipv6 and drop through to ipv4 if poss
			if ipadd.IP.To4() != nil {
				IP = ipadd.IP.String()
			}
		}
	}

	return IP, nil

}

func md5String(input string) string {
	hash := md5.New()
	io.WriteString(hash, input)
	return hex.EncodeToString((hash.Sum(nil)))
}
