Mothership
=========

## About

Mothership is a process monitoring system.


## Installation

Server: 
 1) go run server.go
   

Client:
 1) go build client.go 
 2) Copy binary to clients
 3) run `client {command} {args}`


## To do

 [ ] - Add config file and config parser
 [x] - Gracefully handle failed connections
 [ ] - Add return codes and messages
 [ ] - Persist process state
 [ ] - Move processes to "finished"


