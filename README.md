# Beacon (poc)

Very simple and naive Proof Of Concept to verify the behavior of the method `navigator.sendBeacon()` to track visits.
Checks the both sides, server and client by trying to call JavaScript stuff on the page called by this method.
As expected, on the page called by this method, only the code on server side works.  


## Installation

```
$ go get -u github.com/rvflash/beacon
```

## Quick start

```
$ cd $GOPATH/src/github.com/rvflash/beacon
$ go build && BEACON_PORT=8080 ./beacon
```

As you can see, you can change the port on the server by using the environment variable named BEACON_PORT.
By default, the server is launched on localhost:8080.