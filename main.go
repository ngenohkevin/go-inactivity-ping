package main

import (
	"github.com/ngenohkevin/go-inactivity-ping/ping"
	"log"
	"time"
)

func main() {
	// Create a ticker that triggers every 10 minutes.
	ticker := time.NewTicker(15 * time.Minute)

	// Create a channel to receive results.
	results := make(chan ping.Result)

	list := []string{
		"https://paybutton.onrender.com/",
		"https://www.arnoderrymovers.co.ke",
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
				// Ping the servers again after 10 minutes.
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
				log.Printf("Status: %s, Url: %-20s, Latency: %s", r.StatusCode, r.URL, r.Latency)
			}
		}
	}
}
