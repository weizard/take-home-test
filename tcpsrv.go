package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var limit = make(chan time.Time, 1)
var queryRateWatcher rateState

// tcpsrv current state
var (
	countLock      sync.RWMutex
	currentRequest int
	processedCount int
)

func init() {
	for i := 0; i < 1; i++ {
		limit <- time.Now()
	}
	queryRateWatcher = rateState{
		lastFlush: time.Now(),
		period:    1 * time.Second,
	}
	go func() {
		for t := range time.Tick(time.Duration(1000/float64(RequestRate)) * time.Millisecond) {
			if len(limit) == 0 {
				limit <- t
			}
		}
	}()
}

func queryExternalAPI(s string) error {
	<-limit
	queryRateWatcher.ticRate(1)
	req, err := http.NewRequest(http.MethodGet, apiEndPoint, nil)
	if err != nil {
		fmt.Printf("[queryExternalAPI][reqError] %s\n", err.Error())
		return err
	}
	param := req.URL.Query()
	param.Add("key", s)
	req.URL.RawQuery = param.Encode()
	cli := &http.Client{}
	resp, respErr := cli.Do(req)
	if respErr != nil {
		fmt.Printf("[queryExternalAPI][respError] %s\n", err.Error())
		return respErr
	} else if resp.StatusCode != 200 {
		fmt.Printf("[queryExternalAPI][%d] %s\n", resp.StatusCode, err.Error())
	}
	return nil
}

func connHandler(c net.Conn) {
	fmt.Printf("[connHandler] conn form: %s\n", c.RemoteAddr())
	bufReader := bufio.NewReader(c)
	defer c.Close()
	for {
		c.SetDeadline(time.Now().Add(RequestTimeout))
		stringData, err := bufReader.ReadString('\n')
		if neterr, ok := err.(net.Error); (ok && neterr.Timeout()) || err == io.EOF {
			fmt.Printf("close conn\n")
			countLock.Lock()
			currentRequest--
			processedCount++
			countLock.Unlock()
			c.Close()
			break
		} else if err != nil {
			log.Printf("[connHandler] %s\n", err.Error())
		}
		stringData = strings.Trim(stringData, "\n")
		// log.Printf("%s\n", stringData)
		go queryExternalAPI(stringData)
		if string(stringData) == "quit" {
			fmt.Printf("close conn\n")
			countLock.Lock()
			currentRequest--
			processedCount++
			countLock.Unlock()
			c.Close()
			break
		}
	}
}

// TCPSrv tcp server
func TCPSrv() error {
	s, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		fmt.Printf("tcp server create failed!\n")
		fmt.Printf("Err: %s\n", err.Error())
		return err
	}
	fmt.Printf("server start\n")

	for {
		conn, err := s.Accept()
		countLock.Lock()
		currentRequest++
		countLock.Unlock()
		if err != nil {
			fmt.Printf("Conn error!\n")
			fmt.Printf("Err: [%s] %s\n", conn.RemoteAddr(), err.Error())
			continue
		}
		go connHandler(conn)
	}
}
