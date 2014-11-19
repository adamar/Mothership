
package main


import (
	"net/http"
        "log"
        "io/ioutil"
        "github.com/unrolled/render"
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


func handleEnd(w http.ResponseWriter, req *http.Request) {
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



func mainHandler(w http.ResponseWriter, req *http.Request) {

        r := render.New(render.Options{})
        r.HTML(w, http.StatusOK, "main", nil)

}


func main() {
        http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/heartbeat", handleHeartbeat)
	http.HandleFunc("/end", handleEnd)
        http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", nil)
}





