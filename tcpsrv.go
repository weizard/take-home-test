package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
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
	processedJobs  int
	remainingJobs  int
	failedJobs     int
)

func init() {
	// init limit
	for i := 0; i < 1; i++ {
		limit <- time.Now()
	}
	// Watcher for rps of quering external API
	queryRateWatcher = rateState{
		lastFlush: time.Now(),
		period:    1 * time.Second,
	}
	// control rps
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
		countLock.Lock()
		failedJobs++
		countLock.Unlock()
		fmt.Printf("[queryExternalAPI][reqError] %s\n", err.Error())
		errorLog("queryExternalAPI", "reqError: "+err.Error())
		return err
	}
	param := req.URL.Query()
	param.Add("key", s)
	req.URL.RawQuery = param.Encode()
	cli := &http.Client{}
	resp, cliErr := cli.Do(req)
	if cliErr != nil || resp.StatusCode != 200 {
		countLock.Lock()
		failedJobs++
		remainingJobs--
		countLock.Unlock()
		if cliErr != nil {
			fmt.Printf("[queryExternalAPI][cliError] %s\n", cliErr.Error())
			errorLog("queryExternalAPI", "cliError: "+cliErr.Error())
			return cliErr
		}
		fmt.Printf("[queryExternalAPI][%d] %s\n", resp.StatusCode, err.Error())
		errorLog("respErr", strconv.Itoa(resp.StatusCode))
		return errors.New(string(resp.StatusCode))
	}
	countLock.Lock()
	remainingJobs--
	processedJobs++
	countLock.Unlock()
	return nil
}

func connHandler(c net.Conn) {
	fmt.Printf("[connHandler] conn form: %s\n", c.RemoteAddr())
	bufReader := bufio.NewReader(c)
	defer c.Close()
	for {
		c.SetDeadline(time.Now().Add(RequestTimeout))
		stringData, err := bufReader.ReadString('\n')
		stringData = strings.Trim(stringData, "\n")
		fmt.Println(time.Now(), " ", stringData)
		if neterr, ok := err.(net.Error); (ok && neterr.Timeout()) || err == io.EOF || stringData == "quit" {
			fmt.Printf("%v ", time.Now())
			fmt.Printf("close conn\n")
			c.Close()
			countLock.Lock()
			currentRequest--
			processedCount++
			countLock.Unlock()
			break
		} else if err != nil {
			log.Printf("[connHandler] %s\n", err.Error())
			errorLog("connHandler", err.Error())
		}
		countLock.Lock()
		remainingJobs++
		countLock.Unlock()
		go queryExternalAPI(stringData)
	}
}

// TCPSrv tcp server
func TCPSrv() error {
	s, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		fmt.Printf("tcp server create failed!\n")
		fmt.Printf("Err: %s\n", err.Error())
		errorLog("TCPSrv", "TCP Server Start Error: "+err.Error())
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

// GetCurrentRequest get currentRequest
func GetCurrentRequest() int {
	countLock.RLock()
	defer countLock.RUnlock()
	return currentRequest
}

// GetProcessedCount get processedCount
func GetProcessedCount() int {
	countLock.RLock()
	defer countLock.RUnlock()
	return processedCount
}

// GetProcessedJobs get processedJobs
func GetProcessedJobs() int {
	countLock.RLock()
	defer countLock.RUnlock()
	return processedJobs
}

// GetRemainingJobs get remainingJobs
func GetRemainingJobs() int {
	countLock.RLock()
	defer countLock.RUnlock()
	return remainingJobs
}

// GetFailedJobs get failedJobs
func GetFailedJobs() int {
	countLock.RLock()
	defer countLock.RUnlock()
	return failedJobs
}
