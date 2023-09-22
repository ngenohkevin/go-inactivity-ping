package ping

import (
	"net/http"
	"time"
)

type Result struct {
	URL        string
	Err        error
	Latency    time.Duration
	StatusCode string
	Timestamp  time.Time
}

func Server(url string, ch chan<- Result) {
	start := time.Now()

	if resp, err := http.Get(url); err != nil {
		ch <- Result{
			URL:     url,
			Err:     err,
			Latency: 0,
			//StatusCode: resp.Status,
		}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		ch <- Result{
			URL:        url,
			Err:        nil,
			Latency:    t,
			StatusCode: resp.Status,
		}
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}
}
