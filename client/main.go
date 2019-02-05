/*
TCP Client
*/
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	conn, err := net.Dial("tcp", "localhost:9527")
	if err != nil {
		fmt.Printf("connect to server fail.\n")
		fmt.Printf("Err: %s\n", err.Error())
	}
	sendData := ""
	time.Sleep(2 * time.Second)
	for {
		data := rand.Int()
		fmt.Printf("Sent: %d\n", data)
		sendData += strconv.Itoa(data) + "\n"
		// fmt.Fprintf(conn, strconv.Itoa(data)+"\n")
		conn.Write([]byte(strconv.Itoa(data) + "\n"))

		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
		if data%100 == 0 {
			// conn.Write([]byte(sendData + "quit"))
			// conn.Write([]byte("quit\n"))
			// fmt.Printf("Sent: quit\n")
			break
		}
	}
	sleepTime := 1
	for {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		log.Printf("Sent: test\n")
		_, err = conn.Write([]byte("test\n"))
		if err != nil {
			fmt.Println(err)
			break
		}
		sleepTime++
	}
	// time.Sleep(15 * time.Second)
	// fmt.Printf("Sent: test\n")
	// _, err = conn.Write([]byte("test\n"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// time.Sleep(15 * time.Second)
	// fmt.Printf("Sent: test\n")
	// _, err = conn.Write([]byte("test\n"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// time.Sleep(15 * time.Second)
	// fmt.Printf("Sent: test\n")
	// _, err = conn.Write([]byte("test\n"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// conn.Close()
}
