package main

import (
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestQueryExternalAPI(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	stringData := strconv.Itoa(rand.Int())
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", apiEndPoint, func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(200, "")
		if req.URL.Query().Get("key") != stringData {
			t.Error("wrong data input!")
		}
		if err != nil {
			t.Error(err)
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})
	if err := queryExternalAPI(stringData); err != nil {
		t.Error(err)
	}

}

func TestRPSQueryExternalAPI(t *testing.T) {
	wg := sync.WaitGroup{}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", apiEndPoint,
		func(req *http.Request) (*http.Response, error) {
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
		go queryExternalAPI(sendData)
	}
	wg.Wait()
	rate := float64(reqs) / time.Since(startTime).Seconds()

	if math.Floor(rate) > 30 {
		t.Errorf("rate exceeded rate limit. reqs: %d, rate: %f", reqs, rate)
	} else {
		t.Logf("reqs: %d, rate: %f", reqs, rate)
	}

}
