package main

import (
	"bytes"
	"encoding/json"
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

                startCom()
		args := []string(os.Args[2:])
		data, err := exec.Command(os.Args[1], args...).CombinedOutput()
                log.Print(err)
                log.Print(string(data))
	}

}




func genUuid() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}
        return string(out), nil
}


func startCom() error {

        uuid, _ := genUuid()

        blob := ProcStart{
                UUID: uuid,
                LocalTime: "12:02",
                Command: "/bin/bash /var/www/thing.html",
                Hostname: "ubuntu-server",
                IPaddress: "142.32.12.122",
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


func postJSON(endpoint string, jsonBlob []byte) {

	url := "http://localhost:8080" + endpoint
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBlob))

	req.Header.Set("X-Custom-Header", "MyServer")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

}
