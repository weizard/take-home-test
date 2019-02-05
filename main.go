package main

import (
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

// Config
var (
	RequestRate    = 30               // request per second
	RequestTimeout = 10 * time.Second // second
	Host           = ""               // ip addr
	Port           = "9527"           // port
)

var apiEndPoint = "http://localhost:8080/echo"

func main() {
	var g errgroup.Group
	g.Go(func() error {
		return HTTPSrv()
	})
	g.Go(func() error {
		return TCPSrv()
	})
	if err := g.Wait(); err != nil {
		log.Printf("%s\n", err)
	}

}
