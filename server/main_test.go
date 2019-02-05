package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestQueryExternalAPI(t *testing.T) {
	for rounds := 0; rounds < 1; rounds++ {
		wg := sync.WaitGroup{}
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		count := 0
		go func() {
			for t := range time.Tick(1 * time.Second) {
				fmt.Println(t, ":", count)
				count = 0
			}
		}()
		tempString := ""
		httpmock.RegisterResponder("GET", "http://api.test.com",
			func(req *http.Request) (*http.Response, error) {
				tempString += fmt.Sprintln(time.Now())
				count++
				resp, err := httpmock.NewJsonResponse(200, "")
				if err != nil {
					wg.Done()
					t.Error("500")
					return httpmock.NewStringResponse(500, ""), nil
				}
				wg.Done()
				return resp, nil
			},
		)

		rand.Seed(time.Now().UnixNano())
		startTime := time.Now()
		reqs := rand.Intn(300)
		for i := 0; i < reqs; i++ {
			wg.Add(1)
			data := rand.Int()
			sendData := strconv.Itoa(data)
			go QueryExternalAPI(sendData)
		}
		wg.Wait()
		rate := float64(reqs) / time.Since(startTime).Seconds()

		if math.Floor(rate) > 30 {
			t.Errorf("rate exceeded rate limit. reqs: %d, rate: %f", reqs, rate)
			f, _ := os.Create("log")
			f.WriteString(tempString)
			f.Close()
		} else {
			t.Logf("reqs: %d, rate: %f", reqs, rate)
		}

	}
}
