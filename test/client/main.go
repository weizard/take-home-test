/*
TCP Client
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

func client(wg *sync.WaitGroup) {
	conn, err := net.Dial("tcp", "localhost:9527")
	if err != nil {
		fmt.Printf("connect to server fail.\n")
		fmt.Printf("Err: %s\n", err.Error())
	}
	// sendData := ""
	time.Sleep(2 * time.Second)
	for i := 0; i < 100; i++ {
		data := rand.Int()
		fmt.Printf("Sent: %d\n", data)
		// sendData += strconv.Itoa(data) + "\n"
		// fmt.Fprintf(conn, strconv.Itoa(data)+"\n")
		conn.Write([]byte(strconv.Itoa(data) + "\n"))

		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
		if data%100 == 0 {
			conn.Write([]byte("quit\n"))
			fmt.Printf("Sent: quit\n")
			break
		}
	}
	sleepTime := 1000
	for {
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		log.Printf("Sent: test\n")
		_, err = conn.Write([]byte("test\n"))
		if err != nil {
			fmt.Println(err)
			break
		}
		sleepTime += 1500
	}
	wg.Done()
}

func main() {
	clientNum := flag.Int("n", 1, "number of client you want")
	flag.Parse()
	wg := sync.WaitGroup{}
	rand.Seed(time.Now().UnixNano())
	for n := 0; n < *clientNum; n++ {
		wg.Add(1)
		go client(&wg)
	}
	wg.Wait()
}
