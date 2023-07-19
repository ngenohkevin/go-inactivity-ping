package ping

import (
	"net/http"
	"time"
)

type Result struct {
	URL        string
	Err        error
	Latency    time.Duration
	StatusCode int
	TimeStamp  time.Time
}

func PayBPing(url string, ch chan<- Result) {
	attempts := 3

	start := time.Now()

	if resp, err := http.Get(url); err != nil {
		ch <- Result{
			URL:        url,
			Err:        err,
			Latency:    0,
			StatusCode: http.StatusInternalServerError,
			TimeStamp:  time.Now(),
		}
	} else {
		t := time.Since((start).Round(time.Millisecond))
		ch <- Result{
			URL:        url,
			Err:        checkStatusCode(resp.StatusCode),
			Latency:    t,
			StatusCode: resp.StatusCode,
			TimeStamp:  time.Now(),
		}
		err = resp.Body.Close()
	}

}
func checkStatusCode(statusCode int) error {
	if statusCode == http.StatusServiceUnavailable {
		return http.ErrServerClosed
	}
	return nil
}
