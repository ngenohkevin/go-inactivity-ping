package main

import (
	"github.com/ngenohkevin/go-inactivity-ping/ping"
	"log"
	"time"
)

//const url = "https://paybutton.onrendevr.com/"

func main() {

	stopper := time.After(10 * time.Minute)

	results := make(chan ping.Result)
	list := []string{
		"https://paybutton.onrender.com/",
	}

	for _, url := range list {
		go ping.Server(url, results)
	}
	for range list {
		select {
		case r := <-results:
			if r.Err != nil {
				log.Printf("%s", r.Err)
			} else {
				log.Printf("Status: %s, Url: %-20s, Latency: %s", r.StatusCode, r.URL, r.Latency)
			}
		case t := <-stopper:
			log.Fatalf("timeout %s", t)
		}
	}
}
