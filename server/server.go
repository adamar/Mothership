package main


import (
	"net/http"
        "log"
        "io/ioutil"
        "github.com/unrolled/render"
)


type Broker struct {
    clients         map[chan string]bool
    newClients      chan chan string
    defunctClients  chan chan string
    messages        chan string
}


var broker *Broker = NewBroker()


func (b *Broker) Start() {
    go func() {
        for {
            select {
            case s := <-b.newClients:
                b.clients[s] = true
                log.Println("Added new client")

            case s := <-b.defunctClients:
                delete(b.clients, s)
                log.Println("Removed client")

            case msg := <-b.messages:
                for s, _ := range b.clients {
                    s <- msg
                }
                log.Printf("Broadcast message to %d clients", len(b.clients))
	    }
        }
    }()
}




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
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/heartbeat", handleHeartbeat)
	http.HandleFunc("/end", handleEnd)
        http.HandleFunc("/", mainHandler)
	http.ListenAndServe(":8080", nil)
}


