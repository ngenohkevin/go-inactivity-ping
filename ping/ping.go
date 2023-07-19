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

const maxAttempts = 5                  // Maximum number of retry attempts
const initialBackoff = 1 * time.Second // Initial waiting time before the first retry

func PayBPing(url string, ch chan<- Result) {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		start := time.Now()

		resp, err := http.Get(url)
		if err != nil || resp.StatusCode == http.StatusServiceUnavailable {
			if shouldRetry(attempt) {
				backoff := backoffDuration(attempt)
				time.Sleep(backoff)
				continue
			}

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
				Err:        checkStatusCode(resp.StatusCode),
				Latency:    t,
				StatusCode: resp.StatusCode,
				TimeStamp:  time.Now(),
			}
			err = resp.Body.Close()
			return
		}

		if shouldRetry(attempt) {
			time.Sleep(backoffDuration(attempt))
			continue
		}

		ch <- Result{
			URL:        url,
			Err:        checkStatusCode(resp.StatusCode),
			Latency:    t,
			StatusCode: resp.StatusCode,
			TimeStamp:  time.Now(),
		}

		err = resp.Body.Close()
		return
	}
}

func shouldRetry(attempt int) bool {
	return attempt < maxAttempts-1
}

func backoffDuration(attempt int) time.Duration {
	// Exponential backoff formula: initialBackoff * 2^(attempt-1)
	return initialBackoff << uint(attempt)
}

func checkStatusCode(statusCode int) error {
	if statusCode == http.StatusServiceUnavailable {
		return http.ErrServerClosed
	}
	return nil
}
