package ping

import (
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	URL        string
	Err        error
	Latency    time.Duration
	StatusCode string
	Timestamp  time.Time
}

func Server(url string, ch chan<- Result) {
	start := time.Now()

	// Get timeout from environment variable or use default
	timeout := getEnvTimeout("HTTP_TIMEOUT", 10)
	torTimeout := getEnvTimeout("TOR_TIMEOUT", 60) // Tor needs more time

	//Check if the URL is a .onion URL
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	if isOnionURL(url) {
		// Get Tor proxy address from environment or use default
		torProxy := os.Getenv("TOR_PROXY")
		if torProxy == "" {
			torProxy = "127.0.0.1:9050"
		}

		log.Printf("Connecting to %s via Tor proxy %s with timeout %d seconds", url, torProxy, torTimeout)

		// Test Tor connectivity before proceeding
		if err := testTorConnectivity(torProxy); err != nil {
			ch <- Result{
				URL:       url,
				Err:       fmt.Errorf("tor SOCKS proxy not available at %s: %v", torProxy, err),
				Latency:   0,
				Timestamp: time.Now(),
			}
			return
		}

		// Create a SOCKS5 dialer for Tor
		dialer, err := proxy.SOCKS5("tcp", torProxy, nil, proxy.Direct)
		if err != nil {
			ch <- Result{
				URL:       url,
				Err:       fmt.Errorf("error creating SOCKS5 dialer: %v", err),
				Latency:   0,
				Timestamp: time.Now(),
			}
			return
		}

		// Create context with timeout
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}

		// Use a custom transport with longer timeouts for Tor
		httpTransport := &http.Transport{
			DialContext:           dialContext,
			ResponseHeaderTimeout: time.Duration(torTimeout/2) * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			IdleConnTimeout:       time.Duration(torTimeout) * time.Second,
			TLSHandshakeTimeout:   time.Duration(torTimeout/3) * time.Second,
		}

		client = &http.Client{
			Transport: httpTransport,
			Timeout:   time.Duration(torTimeout) * time.Second, // Longer timeout for Tor
		}
	}

	// Make the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- Result{
			URL:       url,
			Err:       fmt.Errorf("error creating request: %v", err),
			Latency:   0,
			Timestamp: time.Now(),
		}
		return
	}

	// Add a user agent
	req.Header.Set("User-Agent", "Go-Inactivity-Ping/1.0")

	resp, err := client.Do(req)
	if err != nil {
		ch <- Result{
			URL:       url,
			Err:       err,
			Latency:   0,
			Timestamp: time.Now(),
		}
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body for %s: %v", url, err)
		}
	}(resp.Body)

	// Calculate latency and return the result
	latency := time.Since(start).Round(time.Millisecond)
	ch <- Result{
		URL:        url,
		Err:        nil,
		Latency:    latency,
		StatusCode: resp.Status,
		Timestamp:  time.Now(),
	}
}

// testTorConnectivity attempts to connect to the Tor SOCKS proxy to verify it's running
func testTorConnectivity(torProxy string) error {
	conn, err := net.DialTimeout("tcp", torProxy, 5*time.Second)
	if err != nil {
		return err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	return nil
}

func isOnionURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return strings.HasSuffix(u.Host, ".onion")
}

// getEnvTimeout gets timeout from environment or returns default
func getEnvTimeout(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}
