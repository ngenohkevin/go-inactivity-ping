package main

import (
	"fmt"
	"github.com/ngenohkevin/go-inactivity-ping/ping"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const url = "https://paybutton.onrender.com/"

func main() {

	results := make(chan ping.Result)
	done := make(chan struct{})

	go func() {
		for {
			ping.PayBPing(url, results)
			time.Sleep(1 * time.Minute)
		}
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

		<-signals

		done <- struct{}{}
	}()

	for {
		select {
		case result := <-results:
			if result.Err != nil {
				fmt.Printf("Error pinging %s at %s: %v\n", result.URL, result.TimeStamp.Format("2006-01-02 15:04:05"), result.Err)
			} else {
				fmt.Printf("Ping to server at %s successful! latency %v, Status Code: %d\n", result.TimeStamp.Format("2006-01-02 15:04:05"), result.Latency, result.StatusCode)
			}
		case <-done:
			fmt.Println("Terminating...")
			return
		}
	}
}
