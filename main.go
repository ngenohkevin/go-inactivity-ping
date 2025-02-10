package main

import (
	"github.com/ngenohkevin/go-inactivity-ping/ping"
	"log"
	"time"
)

func main() {
	// Create a ticker that triggers every 12 minutes.
	ticker := time.NewTicker(10 * time.Minute)

	// Create a channel to receive results.
	results := make(chan ping.Result)

	list := []string{
		"https://paybutton.onrender.com/",
		"http://someonionurl.onion/",
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
				log.Printf("%s", r.Err)
			} else {
				log.Printf("Status: %s, URL: %-20s, Latency: %s", r.StatusCode, r.URL, r.Latency)
			}
		}
	}
}
