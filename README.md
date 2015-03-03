Mothership
=========

## About

Mothership is a process / micro-service monitoring system. The client wraps a process on a remote system and sends info to the Mothership. 


## Installation

#### Server: 
```sh
go run server.go
```   

#### Client:


```sh
go build client.go 
```
Build Client

```sh
scp client server:/usr/local/bin
```
Copy client to

```sh
client {command} {args}
```
Run client 


## To do

 [x] - Add config file and config parser
 [x] - Gracefully handle failed connections
 [ ] - Add return codes and messages
 [x] - Persist process state
 [x] - Move processes to "finished"
 [ ] - Add iterable times to each command for filtering
 [ ] - Catch crtl-c in client
 [ ] - Remove processes after a variable amount of time
 [ ] - Move all messages to standard fields (ie. create a single struct for all client data)
 [ ] - Order events by date started
 [x] - Add hash to process list
 [x] - Truncate UUID on page


