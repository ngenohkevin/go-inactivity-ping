package ping

import (
	"fmt"
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
	attempts := 3 // Number of retry attempts

	for i := 0; i < attempts; i++ {
		start := time.Now()

		resp, err := http.Get(url)
		if err != nil {
			ch <- Result{
				URL:        url,
				Err:        err,
				Latency:    0,
				StatusCode: http.StatusInternalServerError,
				TimeStamp:  time.Now(),
			}
			return
		}

		t := time.Since(start).Round(time.Millisecond)

		if resp.StatusCode == http.StatusOK {
			ch <- Result{
				URL:        url,
				Err:        nil,
				Latency:    t,
				StatusCode: resp.StatusCode,
				TimeStamp:  time.Now(),
			}
			err = resp.Body.Close()
			return
		}

		// Retry after a short delay
		time.Sleep(1 * time.Second)
	}

	ch <- Result{
		URL:        url,
		Err:        fmt.Errorf("max retry attempts reached"),
		Latency:    0,
		StatusCode: http.StatusServiceUnavailable,
		TimeStamp:  time.Now(),
	}
}

//func checkStatusCode(statusCode int) error {
//	if statusCode == http.StatusServiceUnavailable {
//		return http.ErrServerClosed
//	}
//	return nil
//}
