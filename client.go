package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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
    Ping         string `json:"Ping"`
}


type ProcEnd struct {
    UUID         string `json:"uuid"`
    Error        bool   `json:"error"`
    ExitMessage  string `json:"exitmessage"`
}


func main() {

	switch {
	case len(os.Args) == 1:
		log.Print("nothing to run")

	case len(os.Args) == 2:
		log.Print("run single command without args or flags")

		data, err := exec.Command(os.Args[1]).CombinedOutput()
                log.Print(err)
                log.Print(string(data))

	case len(os.Args) > 2:

                sendJSON()
		args := []string(os.Args[2:])
		data, err := exec.Command(os.Args[1], args...).CombinedOutput()
                log.Print(err)
                log.Print(string(data))
	}

}




func genUuid() {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)
}




func sendJSON() {

        blob := ProcStart{
                UUID: "asxasxasx",
                LocalTime: "12:02",
                Command: "/bin/bash /var/www/thing.html",
                Hostname: "ubuntu-server",
                IPaddress: "142.32.12.122",
                Hash: "A736BC202EC3C",
        }

	url := "http://localhost:8080/start"

	jsonBlob, _ := json.Marshal(blob)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBlob))

	req.Header.Set("X-Custom-Header", "MyServer")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

}
