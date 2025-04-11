package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
}

// InitializeFromEnv loads token and chat ID from environment variables
func InitializeFromEnv() error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHAT_ID")

	if token == "" || chatID == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN or TELEGRAM_CHAT_ID not set in environment")
	}

	Initialize(token, chatID)
	return nil
}

// SendMessage sends a message to the configured Telegram chat
func SendMessage(message string) error {
	if !initialized {
		return fmt.Errorf("telegram bot not initialized")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.Token)

	payload := map[string]string{
		"chat_id": config.ChatID,
		"text":    message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling telegram payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending telegram message: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error, status: %d", resp.StatusCode)
	}

	return nil
}

// NotifyError sends a notification about a URL that is down
func NotifyError(url string, errorMsg string) error {
	statusMap.mu.Lock()
	isCurrentlyDown := statusMap.status[url]
	statusMap.mu.Unlock()

	// If already marked as down, don't send another notification
	if isCurrentlyDown {
		return nil
	}

	// Mark as down
	statusMap.mu.Lock()
	statusMap.status[url] = true
	statusMap.mu.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	message := fmt.Sprintf("🔴 ALERT: Service is DOWN\n\nURL: %s\nTime: %s\nError: %s", url, timestamp, errorMsg)

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

	// Mark as up
	statusMap.mu.Lock()
	statusMap.status[url] = false
	statusMap.mu.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	message := fmt.Sprintf("✅ RECOVERED: Service is back online\n\nURL: %s\nTime: %s\nStatus: %s\nLatency: %s",
		url, timestamp, statusCode, latency)

	return SendMessage(message)
}

// GetStatus returns true if a URL is currently marked as down
func GetStatus(url string) bool {
	statusMap.mu.Lock()
	defer statusMap.mu.Unlock()
	return statusMap.status[url]
}
