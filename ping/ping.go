package ping

import (
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
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

	//Check if the URL is a .onion URL
	client := http.DefaultClient
	if isOnionURL(url) {
		torProxy := "127.0.0.1:9050"
		dialer, err := proxy.SOCKS5("tcp", torProxy, nil, proxy.Direct)
		if err != nil {
			ch <- Result{
				URL:     url,
				Err:     fmt.Errorf("error creating SOCKS5 dialer %v", err),
				Latency: 0,
			}
			return
		}
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
		httpTransport := &http.Transport{
			DialContext: dialContext,
		}
		client = &http.Client{
			Transport: httpTransport,
			Timeout:   10 * time.Second,
		}
	}

	resp, err := client.Get(url)
	if err != nil {
		ch <- Result{
			URL:     url,
			Err:     err,
			Latency: 0,
		}
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ch <- Result{
				URL:     url,
				Err:     fmt.Errorf("error closing response body %v", err),
				Latency: 0,
			}
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

func isOnionURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return strings.HasSuffix(u.Host, ".onion")
}
