package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/ngenohkevin/go-inactivity-ping/logging"
	"github.com/ngenohkevin/go-inactivity-ping/ping"
	"github.com/ngenohkevin/go-inactivity-ping/telegram"
	"log"
	"os"
	"strings"
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
	// Set up logging with timestamps
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Starting Go Inactivity Ping")

	// First, try to load from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, will use embedded config or environment variables")
	}

	// Initialize Telegram bot
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	telegramEnabled := true
	if botToken == "" {
		log.Println("Warning: TELEGRAM_BOT_TOKEN not set, Telegram notifications will not be sent")
		telegramEnabled = false
	} else if chatID == "" {
		log.Println("Warning: TELEGRAM_CHAT_ID not set, Telegram notifications will not be sent")
		telegramEnabled = false
	} else {
		// Initialize Telegram bot
		telegram.Initialize(botToken, chatID)
		log.Println("Telegram bot initialized successfully")

		// Send startup notification
		err := telegram.SendMessage("🚀 <b>URL Monitor Bot started</b>\nMonitoring URLs for availability...")
		if err != nil {
			log.Printf("Failed to send startup notification: %v", err)
		}
	}

	// Load the location for the logger.
	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		// If error, fallback to UTC
		loc = time.UTC
		log.Printf("Failed to load timezone, falling back to UTC: %v", err)
	}

	// Get ping interval from environment or use default
	intervalStr := os.Getenv("PING_INTERVAL")
	interval := 10 * time.Minute
	if intervalStr != "" {
		if intervalInt, err := time.ParseDuration(intervalStr); err == nil {
			interval = intervalInt
		}
	}
	log.Printf("Using ping interval: %s", interval)

	// Create a ticker that triggers based on the interval
	ticker := time.NewTicker(interval)

	// Create a channel to receive results
	results := make(chan ping.Result, 10) // Buffer the channel to avoid blocking

	// Get URLs from environment or use defaults
	var list []string
	urlsEnv := os.Getenv("MONITOR_URLS")
	if urlsEnv != "" {
		list = strings.Split(urlsEnv, ",")
		// Trim spaces
		for i, url := range list {
			list[i] = strings.TrimSpace(url)
		}
	} else {
		// Default URLs
		list = []string{
			"https://google.com", // Test with a regular
		}

		// Check if Tor is explicitly enabled
		torEnabled := os.Getenv("ENABLE_TOR")
		if strings.EqualFold(torEnabled, "true") {
			// Add onion sites if Tor is enabled
			list = append(list, []string{}...)
		} else {
			log.Println("Tor monitoring is DISABLED, skipping .onion sites")
		}
	}

	log.Printf("Monitoring %d URLs", len(list))

	// Start the initial pinging of servers.
	for _, url := range list {
		go ping.Server(url, results)
	}

	// Start a goroutine to handle periodic pinging.
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Ticker triggered, pinging all URLs")
				// Ping the servers again after the interval
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
				// Create log entry
				entry := LogEntry{
					Timestamp: time.Now().In(loc).Format(time.RFC3339),
					Level:     "error",
					App:       "ping-monitor",
					URL:       r.URL,
					Error:     r.Err.Error(),
					Message:   "Ping failed",
				}
				jsonMsg, _ := json.Marshal(entry)
				log.Printf("Error pinging %s: %s", r.URL, r.Err.Error())

				// Send logs to Loki
				err := logging.SendLogsToLoki("error", string(jsonMsg))
				if err != nil {
					log.Printf("Failed to send logs to Loki: %v", err)
				}

				// Send Telegram notification if enabled
				if telegramEnabled {
					err := telegram.NotifyError(r.URL, r.Err.Error())
					if err != nil {
						log.Printf("Failed to send Telegram notification: %v", err)
					}
				}
			} else {
				// Create log entry
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
				log.Printf("Successful ping: %s, URL: %-40s, Latency: %s",
					r.StatusCode, fmt.Sprintf("'%s'", r.URL), r.Latency)

				// Send logs to Loki
				err := logging.SendLogsToLoki("info", string(jsonMsg))
				if err != nil {
					log.Printf("Failed to send logs to Loki: %v", err)
				}

				// Send recovery notification if previously down and Telegram is enabled
				if telegramEnabled {
					err := telegram.NotifyRecovery(r.URL, r.StatusCode, r.Latency)
					if err != nil {
						log.Printf("Failed to send Telegram notification: %v", err)
					}
				}
			}
		}
	}
}
