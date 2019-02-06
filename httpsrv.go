package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// CurrentState server's current state
type CurrentState struct {
	CurrentConn   int
	RPS           float64
	ProcessedReq  int
	ProcessedJobs int
	RemainingJobs int
	FailedJobs    int
}

func getState(w http.ResponseWriter, r *http.Request) {
	state := &CurrentState{
		CurrentConn:   GetCurrentRequest(),
		RPS:           queryRateWatcher.getRate(),
		ProcessedReq:  GetProcessedCount(),
		ProcessedJobs: GetProcessedJobs(),
		RemainingJobs: GetRemainingJobs(),
		FailedJobs:    GetFailedJobs(),
	}
	b, _ := json.Marshal(state)
	fmt.Fprintf(w, "%s", b)
}

// HTTPSrv http server for observe tcp server status
func HTTPSrv() error {
	http.HandleFunc("/state", getState)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Printf("http server error: %s\n", err.Error())
	}
	return err
}
