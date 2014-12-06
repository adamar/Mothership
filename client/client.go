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


var UUID = genUuid()


func main() {

	switch {
	case len(os.Args) == 1:
		log.Print("nothing to run")

	case len(os.Args) == 2:

                go sendHeartbeat()
		data, err := exec.Command(os.Args[1]).CombinedOutput()
                log.Print(err)
                log.Print(string(data))

	case len(os.Args) > 2:

                startCom(os.Args[1:])
                go sendHeartbeat()
		args := []string(os.Args[2:])
		data, err := exec.Command(os.Args[1], args...).CombinedOutput()
                log.Print(err)
                log.Print(string(data))
	}

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



func startCom(cmd []string) error {

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

        endpoint := "/start"
        jsonBlob, err := json.Marshal(blob)
        if err != nil {
            return err
        }

        postJSON(endpoint, jsonBlob)
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

        postJSON(endpoint, jsonBlob)

   }

}


func postJSON(endpoint string, jsonBlob []byte) {

	url := "http://localhost:8080" + endpoint
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBlob))

	req.Header.Set("X-Custom-Header", "MyServer")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

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
