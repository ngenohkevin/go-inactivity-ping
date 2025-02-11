package main

import (
	"encoding/json"
	"github.com/ngenohkevin/go-inactivity-ping/logging"
	"github.com/ngenohkevin/go-inactivity-ping/ping"
	"log"
	"time"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	App       string `json:"app"`
	URL       string `json:"url,omitempty"`
	Status    string `json:"status,omitempty"`
	Latency   string `json:"latency,omitempty"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
}

func main() {
	// Load the location for the logger.
	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		// If error, fallback to UTC
		loc = time.UTC
	}

	// Create a ticker that triggers every 10 minutes.
	ticker := time.NewTicker(10 * time.Minute)

	// Create a channel to receive results.
	results := make(chan ping.Result)

	list := []string{
		"https://paybutton.onrender.com/",
		"http://dweb5sm34uzajabrjtl2jvznaaeop526p3sqi73hw4gur5ggkadqa3ad.onion/",
		"http://iqualowxvgqaijkrxy2xrl4peewievvvocenmt2qfbkxcto6cqt2anyd.onion/",
		"http://ajwccz5y5jn33bwt2jiibcytg2pk4gdluwmij5i4qse3oahykq4fj7id.onion/",
	}

	// Start the initial pinging of servers.
	for _, url := range list {
		go ping.Server(url, results)
	}

	// Start a goroutine to handle periodic pinging.
	go func() {
		for {
			select {
			case <-ticker.C:
				// Ping the servers again after 12 minutes.
				for _, url := range list {
					go ping.Server(url, results)
				}
			}
		}
	}()

	// Process results.
	for {
		select {
		case r := <-results:
			if r.Err != nil {
				entry := LogEntry{
					Timestamp: time.Now().In(loc).Format(time.RFC3339),
					Level:     "error",
					App:       "ping-monitor",
					URL:       r.URL,
					Error:     r.Err.Error(),
					Message:   "Ping failed",
				}
				jsonMsg, _ := json.Marshal(entry)
				log.Printf("Error: %s", r.Err.Error())
				_ = logging.SendLogsToLoki("error", string(jsonMsg))
			} else {
				entry := LogEntry{
					Timestamp: time.Now().In(loc).Format(time.RFC3339),
					Level:     "info",
					App:       "ping-monitor",
					URL:       r.URL,
					Status:    r.StatusCode,
					Latency:   r.Latency.String(),
					Message:   "Ping successful",
				}
				jsonMsg, _ := json.Marshal(entry)
				log.Printf("Status: %s, URL: %-20s, Latency: %s", r.StatusCode, r.URL, r.Latency)
				_ = logging.SendLogsToLoki("info", string(jsonMsg))
			}
		}
	}
}
