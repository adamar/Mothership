
package main


import (
	"net/http"
        "log"
        "io/ioutil"
)


func handleStart(w http.ResponseWriter, req *http.Request) {
     if req.Method == "POST" {
         body, err := ioutil.ReadAll(req.Body)
         if err != nil {
            log.Print(err)
         }
         log.Print(string(body))
         w.Write([]byte("status"))
     } else {
         w.Write([]byte("error, do post"))
     }
}

func handleHeartbeat(w http.ResponseWriter, req *http.Request) {
        w.Write([]byte("something"))
}

func handleEnd(w http.ResponseWriter, req *http.Request) {
        w.Write([]byte("something"))
}


func main() {
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/heartbeat", handleHeartbeat)
	http.HandleFunc("/end", handleEnd)
	http.ListenAndServe(":8080", nil)
}





