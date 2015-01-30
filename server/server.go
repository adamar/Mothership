package main


import (
	"net/http"
    "log"
    "fmt"
    "io/ioutil"
    "github.com/unrolled/render"
    "os"
    "encoding/json"
    "github.com/boltdb/bolt"
)


type Broker struct {
    clients         map[chan string]bool
    newClients      chan chan string
    defunctClients  chan chan string
    messages        chan string
}


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



var broker *Broker = NewBroker()
var debug = checkDebugStatus()
var procDB = setupDB()
var procs = []byte("processes")
var defunctprocs = []byte("defunctprocesses")


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

         out, err := unmarshalStart(body)
         if err != nil {
             log.Print(err)
         }

         Put(procs, []byte(out.UUID), body)

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

         out, err := unmarshalEnd(body)
         if err != nil {
             log.Print(err)
         }

         Delete(procs, []byte(out.UUID))
         Put(defunctprocs, []byte(out.UUID), body)


         data := `{"type":"end","body":` + string(body) + `}`
         broker.messages <- data
         w.Write([]byte("status"))
     } else {
         w.Write([]byte("error, do post"))
     }

}


func unmarshalStart(data []byte) (*ProcStart, error) {

    start := &ProcStart{}
    
    if err := json.Unmarshal(data, &start); err != nil {
        return nil, err
    }
    return start, nil

}


func unmarshalEnd(data []byte) (*ProcEnd, error) {

    end := &ProcEnd{}

    if err := json.Unmarshal(data, &end); err != nil {
        return nil, err
    }
    return end, nil

}



func mainHandler(w http.ResponseWriter, req *http.Request) {

    data := GetMany(procs)
    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, "main", data)

}


func defunctHandler(w http.ResponseWriter, req *http.Request) {

    data := GetMany(defunctprocs)
    r := render.New(render.Options{})
    r.HTML(w, http.StatusOK, "defunct", data)

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



func Put(bucket []byte, key []byte, value []byte) error {

    err := procDB.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucket)
        err := b.Put(key, value)
        return err
    }) 

    if err != nil {
        return err
    }

    return nil

}


func Get(bucket []byte, key []byte) []byte {

    var val []byte

    procDB.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucket)
        val = b.Get(key)
        return nil
    })

    return val

}


func Delete(bucket []byte, key []byte) error {

    err := procDB.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucket)
        err := b.Delete(key)
        if err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        return err
    }

    return nil

}




func GetMany(bucket []byte) []ProcStart {

    var data []ProcStart
    var p ProcStart

    procDB.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucket)
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            json.Unmarshal(v, &p)
            data = append(data, p)
        }

        return nil
    })

    return data

}



func GetSince(bucket []byte, t1 time.Time) {

    var data []ProcStart
    var p ProcStart

    procDB.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(bucket)
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            json.Unmarshal(v, &p)
            data = append(data, p)

            t0, _ := time.Parse(time.RFC3339, data.curTime)

            var duration time.Duration = t1.Sub(t0)
        }

        return nil
    })

    return data

}



func setupDB() *bolt.DB {

    DB, err := bolt.Open("bolt.db", 0644, nil)
    if err != nil {
        log.Fatal(err)
    }

    DB.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucket(procs)
        if err != nil {
            return fmt.Errorf("create bucket: %s", err)
        }
        return nil
    })

    DB.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucket(defunctprocs)
        if err != nil {
            return fmt.Errorf("create bucket: %s", err)
        }
        return nil
    })


    return DB

}


func main() {

  
    broker.Start()

    http.Handle("/events/", broker)
    http.HandleFunc("/start", handleStart)
    http.HandleFunc("/heartbeat", handleHeartbeat)
    http.HandleFunc("/end", handleEnd)

    http.HandleFunc("/", mainHandler)
    http.HandleFunc("/defunct", defunctHandler)
 
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    serverAdd := ":8080"
    fmt.Printf("Server listening on %s\n", serverAdd)

    http.ListenAndServe(serverAdd, nil)

}


