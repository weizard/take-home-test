/*
TCP Server
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Config
var (
	RequestRate    = 30     // request per second
	RequestTimeout = 10     // second
	Host           = ""     // ip addr
	Port           = "9527" // port
)

var (
	countLock      sync.RWMutex
	currentRequest int
	processedCount int
)

var rateWatch rateState

type rateState struct {
	mu        sync.RWMutex
	lastFlush time.Time
	period    time.Duration
	tickCount int
	tps       float64
}

func (pRate *rateState) ticRate(i int) {
	pRate.mu.Lock()
	pRate.tickCount += i
	now := time.Now()
	tps := 0.
	if now.Sub(pRate.lastFlush) >= pRate.period {
		if pRate.tickCount > 0 {
			tps = float64(pRate.tickCount) / now.Sub(pRate.lastFlush).Seconds()
		}
		pRate.tickCount = 0
		pRate.lastFlush = now
		pRate.tps = tps
	}
	pRate.mu.Unlock()
}

func (pRate *rateState) getRate() float64 {
	pRate.ticRate(0)
	t := 0.
	pRate.mu.RLock()
	t = pRate.tps
	pRate.mu.RUnlock()
	return t
}

var apiEndPoint = "http://api.test.com"

// ConnHandler process the conn
func ConnHandler(c net.Conn) {
	fmt.Println(c.RemoteAddr())
	bufReader := bufio.NewReader(c)
	for {
		c.SetDeadline(time.Now().Add(5 * time.Second))
		stringData, err := bufReader.ReadString('\n')
		if neterr, ok := err.(net.Error); (ok && neterr.Timeout()) || err == io.EOF {
			fmt.Printf("close conn\n")
			countLock.Lock()
			currentRequest--
			processedCount++
			countLock.Unlock()
			c.Close()
			break
		}
		if err != nil {
			log.Println(err)
		}
		stringData = strings.Trim(stringData, "\n")
		log.Println(stringData)
		go QueryExternalAPI(stringData)
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

var limit = make(chan time.Time, 1)

// QueryExternalAPI query external api
func QueryExternalAPI(s string) {
	<-limit
	rateWatch.ticRate(1)
	req, err := http.NewRequest(http.MethodGet, apiEndPoint, nil)
	if err != nil {
		fmt.Println(err)
	}
	param := req.URL.Query()
	param.Add("key", s)
	req.URL.RawQuery = param.Encode()
	cli := &http.Client{}
	resp, respErr := cli.Do(req)
	// if neterr, ok := err.(net.Error); ok {
	// 	fmt.Println(neterr.Error())
	// 	fmt.Println(err)
	// }
	if respErr != nil {
		fmt.Println("1:", respErr)
	} else if resp.StatusCode != 200 {
		fmt.Println("2", err)
	}
}

func init() {
	for i := 0; i < 1; i++ {
		limit <- time.Now()
	}
	rateWatch = rateState{
		lastFlush: time.Now(),
		period:    1 * time.Second,
	}
	go func() {
		for t := range time.Tick(time.Duration(1000/float64(RequestRate)) * time.Millisecond) {
			if len(limit) == 0 {
				limit <- t
			}
			// rateWatch.ticRate(0)
		}
	}()
}

type Stat struct {
	CurrentConn  int
	RPS          float64
	ProcessedReq int
}

func getStat(w http.ResponseWriter, r *http.Request) {

	rps := rateWatch.getRate()
	countLock.RLock()
	stat := &Stat{
		CurrentConn:  currentRequest,
		RPS:          rps,
		ProcessedReq: processedCount,
	}
	countLock.RUnlock()
	b, _ := json.Marshal(stat)
	fmt.Fprintf(w, "%s", b)
}

func httpSrv() {
	http.HandleFunc("/stat", getStat)
	http.ListenAndServe(":80", nil)
}

func main() {

	go httpSrv()
	s, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		fmt.Printf("tcp server create failed!\n")
		fmt.Printf("Err: %s\n", err.Error())
	}
	fmt.Printf("server start\n")

	for {
		conn, err := s.Accept()
		countLock.Lock()
		currentRequest++
		countLock.Unlock()
		if err != nil {
			fmt.Printf("Conn error!\n")
			fmt.Printf("Err: %s\n", err.Error())
			continue
		}
		go ConnHandler(conn)
	}

}
