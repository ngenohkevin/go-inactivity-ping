# **Go Inactivity Ping**

Go Inactivity Ping is a lightweight, command-line tool written in Go that periodically pings servers to ensure they remain reachable. Designed to prevent server spin-down due to inactivity and monitor uptime, it includes support for `.onion` sites hosted on the Tor network and Telegram notifications.

## **About the Project**

Many hosting services automatically spin down servers after a period of inactivity, which can disrupt applications. Go Inactivity Ping ensures your server remains active by sending periodic ping requests at configurable intervals. Additionally, it offers monitoring for `.onion` (Tor) sites, helping you ensure uptime for your hidden services.

The tool uses Go's native `net/http` library, along with a SOCKS5 dialer for `.onion` support, to make HTTP requests to your server. By default, it sends a ping every **10 minutes**, though this interval can be customized.

---

## **Features**

- **Periodic Ping Requests**: Keeps your server active by sending requests at regular intervals.
- **.onion Site Monitoring**: Supports `.onion` URLs with Tor network compatibility via SOCKS5 proxies.
- **Latency Tracking**: Measures and logs the response latency for each ping.
- **Uptime Monitoring**: Monitors the availability of your server and detects connectivity issues.
- **Response Logging**: Detailed response logs, including timestamps, status codes, and error messages, for analysis.
- **Telegram Notifications**: Receive instant alerts when services go down and recover.
- **Loki Integration**: Sends logs to Grafana Loki for centralized log management.

---

## **Setup**

1. Clone the repository
2. Configure your `.env` file with the required credentials:
   ```
   # Loki Configuration (for logging)
   LOKI_USER=your_loki_user
   LOKI_API_KEY=your_loki_api_key
   LOKI_URL=your_loki_url
   
   # Telegram Configuration (for notifications)
   TELEGRAM_BOT_TOKEN=your_telegram_bot_token
   TELEGRAM_CHAT_ID=your_telegram_chat_id
   
   # Timeout settings (in seconds)
   HTTP_TIMEOUT=10
   TOR_TIMEOUT=60
   
   # Tor proxy configuration
   TOR_PROXY=127.0.0.1:9050
   # Set to "true" to enable monitoring of .onion sites
   ENABLE_TOR=false
   
   # Ping interval (supports Go duration format: 10m, 1h, etc.)
   PING_INTERVAL=10m
   
   # Comma-separated list of URLs to monitor (optional)
   MONITOR_URLS=https://example.com,https://example2.com
   ```

3. Build and run the application:
   ```
   go build
   ./go-inactivity-ping
   ```

## **Tor Setup for .onion Sites**

To monitor `.onion` sites, you need to have Tor installed and running:

1. Install Tor:
   - Ubuntu/Debian: `sudo apt-get install tor`
   - macOS (with Homebrew): `brew install tor`
   - Windows: Download the Tor Browser and use its bundled Tor service

2. Make sure Tor is running and listening on port 9050 (default SOCKS proxy port)

3. In your `.env` file, set `ENABLE_TOR=true`

4. Check that `TOR_PROXY` is set to `127.0.0.1:9050` (or your custom Tor SOCKS proxy address)

5. You may need to increase `TOR_TIMEOUT` value for .onion sites as they can be slower to respond

## **Telegram Bot Setup**

1. Create a Telegram bot using the [BotFather](https://t.me/botfather) on Telegram:
   - Start a chat with BotFather
   - Send `/newbot` command
   - Follow the instructions to create a bot
   - BotFather will give you a token (keep it secret!)

2. Get your bot token and add it to the `.env` file under `TELEGRAM_BOT_TOKEN`

3. To get your chat ID:
   - Start a chat with your bot
   - Send a message to your bot
   - Visit `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
   - Look for the `chat` object and copy the `id` value
   - Add this ID to the `TELEGRAM_CHAT_ID` in your `.env` file

Once configured, you'll receive notifications when:
- Any monitored URL goes down
- A previously down URL comes back online

## **Troubleshooting**

### Tor Connectivity Issues

If you're having problems with `.onion` sites:

1. Make sure Tor is running:
   - Check with `ps -ef | grep tor` 
   - Or check if port 9050 is open: `telnet 127.0.0.1 9050`

2. Verify Tor proxy settings in `.env`

3. Increase timeout values for Tor:
   - Set `TOR_TIMEOUT=120` or higher

4. Test with a known reliable .onion site first

### Telegram Notification Issues

1. Verify your bot token format
   
2. Make sure you've started a conversation with the bot

3. Confirm your chat ID is correctly set

4. Check the application logs for detailed error messages

---

## **Customization**

You can customize the behavior of the tool through environment variables in the `.env` file:

- `PING_INTERVAL`: How often to check sites (e.g., `10m`, `1h`)
- `HTTP_TIMEOUT`: Timeout for regular websites in seconds
- `TOR_TIMEOUT`: Timeout for .onion sites in seconds
- `MONITOR_URLS`: Comma-separated list of URLs to monitor
- `ENABLE_TOR`: Set to `true` to enable .onion site monitoring

---

## **Planned Features**
- **Health Dashboard**: Visualize server uptime, response times, and error patterns using an intuitive dashboard.
- **Retry Mechanism**: Automatically retry failed pings before marking a server as unreachable.
- **Web Interface**: A simple web-based dashboard for monitoring and configuration.

---
