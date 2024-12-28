# **Go Inactivity Ping**

Go Inactivity Ping is a lightweight, command-line tool written in Go that periodically pings servers to ensure they remain reachable. Designed to prevent server spin-down due to inactivity and monitor uptime, it now includes support for `.onion` sites hosted on the Tor network.

## **About the Project**

Many hosting services automatically spin down servers after a period of inactivity, which can disrupt applications. Go Inactivity Ping ensures your server remains active by sending periodic ping requests at configurable intervals. Additionally, it offers monitoring for `.onion` (Tor) sites, helping you ensure uptime for your hidden services.

The tool uses Go's native `net/http` library, along with a SOCKS5 dialer for `.onion` support, to make HTTP requests to your server. By default, it sends a ping every **12 minutes**, though this interval can be customized.

---

## **Features**

- **Periodic Ping Requests**: Keeps your server active by sending requests at regular intervals.
- **.onion Site Monitoring**: Supports `.onion` URLs with Tor network compatibility via SOCKS5 proxies.
- **Latency Tracking**: Measures and logs the response latency for each ping.
- **Uptime Monitoring**: Monitors the availability of your server and detects connectivity issues.

### **Planned Features**
- **Response Logging**: Save detailed response logs, including timestamps, status codes, and error messages, for analysis.
- **Alerting System**: Configure notifications (e.g., email or Slack) to alert you when the server becomes unreachable or latency exceeds a threshold.
- **Health Dashboard**: Visualize server uptime, response times, and error patterns using an intuitive dashboard.
- **Retry Mechanism**: Automatically retry failed pings before marking a server as unreachable.

---

