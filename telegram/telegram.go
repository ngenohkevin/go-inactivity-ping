package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// BotConfig holds the Telegram bot configuration
type BotConfig struct {
	Token  string
	ChatID string
}

// StatusTracker keeps track of URL status to avoid repeated notifications
type StatusTracker struct {
	mu     sync.Mutex
	status map[string]bool // true if URL is down, false if UP
}

var (
	config      BotConfig
	statusMap   = &StatusTracker{status: make(map[string]bool)}
	initialized bool
)

// Initialize sets up the Telegram bot with token and chat ID
func Initialize(token, chatID string) {
	config.Token = token
	config.ChatID = chatID
	initialized = true

	// Log initialization but mask part of the token for security
	maskedToken := token
	if len(token) > 10 {
		maskedToken = token[:8] + "..." + token[len(token)-4:]
	}
	log.Printf("Telegram bot initialized with token %s and chat ID %s", maskedToken, chatID)
}

// SendMessage sends a message to the configured Telegram chat
func SendMessage(message string) error {
	if !initialized {
		return fmt.Errorf("telegram bot not initialized")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.Token)

	payload := map[string]string{
		"chat_id":    config.ChatID,
		"text":       message,
		"parse_mode": "HTML", // Enable HTML formatting
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling telegram payload: %v", err)
	}

	// Create a client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending telegram message: %v", err)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("Telegram message sent successfully: %s", message[:minimum(30, len(message))]+"...")
	return nil
}

// NotifyError sends a notification about a URL that is down
func NotifyError(url string, errorMsg string) error {
	statusMap.mu.Lock()
	isCurrentlyDown := statusMap.status[url]
	statusMap.mu.Unlock()

	// If already marked as down, don't send another notification
	if isCurrentlyDown {
		log.Printf("URL %s is already marked as down, skipping notification", url)
		return nil
	}

	log.Printf("Marking URL %s as DOWN and sending notification", url)

	// Mark as down
	statusMap.mu.Lock()
	statusMap.status[url] = true
	statusMap.mu.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	// Use better formatting for Telegram
	message := fmt.Sprintf("<b>🔴 ALERT: Service is DOWN</b>\n\n<b>URL:</b> %s\n<b>Time:</b> %s\n<b>Error:</b> %s",
		url, timestamp, errorMsg)

	return SendMessage(message)
}

// NotifyRecovery sends a notification about a URL that has recovered
func NotifyRecovery(url string, statusCode string, latency time.Duration) error {
	statusMap.mu.Lock()
	isCurrentlyDown := statusMap.status[url]
	statusMap.mu.Unlock()

	// If not marked as down, don't send recovery notification
	if !isCurrentlyDown {
		return nil
	}

	log.Printf("Marking URL %s as UP and sending recovery notification", url)

	// Mark as up
	statusMap.mu.Lock()
	statusMap.status[url] = false
	statusMap.mu.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	// Use better formatting for Telegram
	message := fmt.Sprintf("<b>✅ RECOVERED: Service is back online</b>\n\n<b>URL:</b> %s\n<b>Time:</b> %s\n<b>Status:</b> %s\n<b>Latency:</b> %s",
		url, timestamp, statusCode, latency)

	return SendMessage(message)
}

// Helper function for min
func minimum(a, b int) int {
	if a < b {
		return a
	}
	return b
}
