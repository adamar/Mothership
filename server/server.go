package main


import (
	"net/http"
    "log"
    "fmt"
    "io/ioutil"
    "github.com/unrolled/render"
    "os"
)


type Broker struct {
    clients         map[chan string]bool
    newClients      chan chan string
    defunctClients  chan chan string
    messages        chan string
}


var broker *Broker = NewBroker()
var debug = checkDebugStatus()


func (b *Broker) Start() {
    go func() {
        for {
            select {
            case s := <-b.newClients:
                b.clients[s] = true
                if debug == true {
                    log.Println("Added new client")
                }

            case s := <-b.defunctClients:
                delete(b.clients, s)
                if debug == true {
                    log.Println("Removed client")
                }

            case msg := <-b.messages:
                for s, _ := range b.clients {
                    s <- msg
                }
                if debug == true {
                    log.Printf("Broadcast message to %d clients", len(b.clients))
                }

	    }
        }
    }()
}

func checkDebugStatus() bool {

    if os.Getenv("DEBUG") == "TRUE" {
        return true
    }

    return false

}


func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    f, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }

    messageChan := make(chan string)

    b.newClients <- messageChan

    defer func() {
        b.defunctClients <- messageChan
    }()

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    for i := 0; i < 10; i++ {
        msg := <-messageChan
        fmt.Fprintf(w, "data: %s\n\n", msg)
        f.Flush()
    }

    log.Println("Finished HTTP request at ", r.URL.Path)
}



func handleStart(w http.ResponseWriter, req *http.Request) {
     if req.Method == "POST" {
         body, err := ioutil.ReadAll(req.Body)
         if err != nil {
            log.Print(err)
         }

         if debug == true {
             log.Print(string(body))
         }

         data := `{"type":"start","body":` + string(body) + `}`
         broker.messages <- data
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
         if debug == true {
             log.Print(string(body))
         }

         data := `{"type":"heartbeat","body":` + string(body) + `}`
         broker.messages <- data
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
         if debug == true {
             log.Print(string(body))
         }

         data := `{"type":"end","body":` + string(body) + `}`
         broker.messages <- data
         w.Write([]byte("status"))
     } else {
         w.Write([]byte("error, do post"))
     }

}



func mainHandler(w http.ResponseWriter, req *http.Request) {

        r := render.New(render.Options{})
        r.HTML(w, http.StatusOK, "main", nil)

}


func NewBroker() *Broker {

    broker := &Broker{
        make(map[chan string]bool),
        make(chan (chan string)),
        make(chan (chan string)),
        make(chan string),
    }

    return broker

}



func main() {

    broker.Start()

    http.Handle("/events/", broker)
    http.HandleFunc("/start", handleStart)
    http.HandleFunc("/heartbeat", handleHeartbeat)
    http.HandleFunc("/end", handleEnd)
    http.HandleFunc("/", mainHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    serverAdd := ":8080"
    fmt.Printf("Server listening on %s\n", serverAdd)

    http.ListenAndServe(serverAdd, nil)

}


