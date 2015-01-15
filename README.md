Mothership
=========

## About

Mothership is a process / micro-service monitoring system. 


## Installation

Server: 
 1) go run server.go
   

Client:
 1) go build client.go 
 2) Copy binary to clients
 3) run `client {command} {args}`


## To do

 [x] - Add config file and config parser
 [x] - Gracefully handle failed connections
 [ ] - Add return codes and messages
 [ ] - Persist process state
 [ ] - Move processes to "finished"
 [ ] - Add iterable times to each command for filtering
 [ ] - Catch crtl-c in client
 [ ] - Remove processes after a variable amount of time


