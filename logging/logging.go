package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"time"
)

// LokiPushData is a struct that holds the logs to be sent to loki
type LokiPushData struct {
	Streams []LokiStream `json:"streams"`
}

// LokiStream is a struct that holds the stream and values of the logs
type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func SendLogsToLoki(level, message string) error {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	//load credentials from environment variables
	lokiUser := os.Getenv("LOKI_USER")
	lokiPassword := os.Getenv("LOKI_API_KEY")
	lokiURL := os.Getenv("LOKI_URL")

	if lokiUser == "" || lokiPassword == "" || lokiURL == "" {
		return fmt.Errorf("loki credentials not set")
	}

	logTime := fmt.Sprintf("%d", time.Now().UnixNano())
	lokiData := LokiPushData{
		Streams: []LokiStream{
			{
				Stream: map[string]string{
					"level": level,
					"app":   "ping-monitor",
					//add more labels if available (e.g., "server": someServerValue)
				},
				Values: [][]string{{logTime, message}},
			},
		},
	}
	//convert logs to json
	jsonData, err := json.Marshal(lokiData)
	if err != nil {
		return fmt.Errorf("error marshalling logs to json %v", err)
	}
	// create http request to send logs to loki
	req, err := http.NewRequest("POST", lokiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating http request %v", err)
	}
	//set headers and basic auth
	req.SetBasicAuth(lokiUser, lokiPassword)
	req.Header.Set("Content-Type", "application/json")

	//send http request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending logs to loki %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	//check if logs were sent successfully
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error sending logs to loki, status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}
