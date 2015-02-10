package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
    "net"
    "strings"
    "time"
    "strconv"
    "io/ioutil"
)


type ProcStart struct {
    UUID         string `json:"uuid"`
    LocalTime    string `json:"localtime"`
    Command      string `json:"command"`
    Hostname     string `json:"hostname"`
    IPaddress    string `json:"ipaddress"`
    Hash         string `json:"hash"`
}


type Heartbeat struct {
    UUID         string `json:"uuid"`
    Ping         string `json:"Ping"`
    RunningTime  string `json:"runningtime"`
}


type ProcEnd struct {
    UUID         string `json:"uuid"`
    Error        bool   `json:"error"`
    ExitMessage  string `json:"exitmessage"`
}


type Config struct{
    Hostname     string `json:"hostname"`
    Port         string `json:"port"`
}



var UUID = genUuid()
var debug = checkDebugStatus()
var conf = parseConfig()


func main() {

    switch {
    case len(os.Args) == 1:
        if debug == true {
            log.Print("nothing to run")
        }

    case len(os.Args) == 2:

        sendCom(os.Args[1], "/start")
        go sendHeartbeat()
        cmd := exec.Command(os.Args[1])

        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Env = os.Environ()

        cmd.Start()

        cmd.Wait()


    case len(os.Args) > 2:

        sendCom(os.Args[1:], "/start")
        go sendHeartbeat()
        args := []string(os.Args[2:])
        cmd := exec.Command(os.Args[1], args...)

        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Env = os.Environ()

        cmd.Start()

        cmd.Wait()

    }

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
        fixed = strings.Replace(string(out), "\n","", -1)

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

        blob := ProcStart{
                UUID: UUID,
                LocalTime: curTime,
                Command: cmdstring,
                Hostname: hostname,
                IPaddress: ip,
                Hash: "A736BC202EC3C",
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

        blob := Heartbeat{
                Ping: string(time.Now().Format(time.RFC3339)),
                UUID: UUID,
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


func sendEnd() {

        blob := ProcEnd{
                UUID: UUID,
                Error: true,
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

func parseConfig() *Config {

    content, err := ioutil.ReadFile("config.json")
    if err!=nil{
        log.Print(err)
    }

    var conf Config
    err=json.Unmarshal(content, &conf)
    if err!=nil{
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
