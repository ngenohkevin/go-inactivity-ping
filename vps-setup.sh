#!/bin/bash
# VPS setup script for ping-monitor
# This script helps set up the Tor service and prepare for ping-monitor installation

# Exit on error
set -e

# Function to display colored output
print_status() {
  echo -e "\e[1;34m[*] $1\e[0m"
}

print_success() {
  echo -e "\e[1;32m[+] $1\e[0m"
}

print_error() {
  echo -e "\e[1;31m[!] $1\e[0m"
}

print_status "VPS setup for ping-monitor with Tor support"
print_status "==============================================="

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  print_error "Please run as root"
  exit 1
fi

# Create installation directory
INSTALL_DIR="/opt/ping-monitor"
print_status "Creating installation directory: $INSTALL_DIR"
mkdir -p $INSTALL_DIR

# Install Tor if requested
read -p "Do you want to install Tor? (y/n): " install_tor
if [[ $install_tor == "y" || $install_tor == "Y" ]]; then
  print_status "Installing Tor..."
  
  # Detect distribution
  if [ -f /etc/debian_version ]; then
    # Debian/Ubuntu
    apt-get update
    apt-get install -y tor
    systemctl enable tor
    systemctl start tor
  elif [ -f /etc/redhat-release ]; then
    # CentOS/RHEL
    yum install -y epel-release
    yum install -y tor
    systemctl enable tor
    systemctl start tor
  else
    print_error "Unsupported distribution. Please install Tor manually."
    exit 1
  fi
  
  print_success "Tor installed and started"
  
  # Verify Tor is working
  if nc -z -w5 127.0.0.1 9050; then
    print_success "Tor SOCKS proxy is available at 127.0.0.1:9050"
  else
    print_error "Tor SOCKS proxy doesn't seem to be running at 127.0.0.1:9050"
    print_error "Please check Tor installation manually"
  fi
fi

# Copy ping-monitor files if they exist in current directory
if [ -f "ping-monitor" ]; then
  print_status "Copying ping-monitor executable..."
  cp ping-monitor $INSTALL_DIR/
  chmod +x $INSTALL_DIR/ping-monitor
  
  if [ -f ".env" ]; then
    print_status "Copying configuration..."
    cp .env $INSTALL_DIR/
  else
    print_status "No .env file found, creating a template..."
    cat > $INSTALL_DIR/.env << EOL
# Configuration for ping-monitor

# Loki Configuration
LOKI_USER=your_loki_user
LOKI_API_KEY=your_loki_api_key
LOKI_URL=your_loki_url

# Telegram configuration
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_CHAT_ID=your_chat_id

# Timeout settings (in seconds)
HTTP_TIMEOUT=10
TOR_TIMEOUT=60

# Tor proxy configuration
TOR_PROXY=127.0.0.1:9050
# Set to "true" to enable monitoring of .onion sites
ENABLE_TOR=true

# Ping interval (supports Go duration format: 10m, 1h, etc.)
PING_INTERVAL=10m

# Comma-separated list of URLs to monitor
MONITOR_URLS=https://paybutton.onrender.com/,https://google.com
EOL
  fi
  
  # Create systemd service
  print_status "Setting up systemd service..."
  cat > /etc/systemd/system/ping-monitor.service << EOL
[Unit]
Description=URL Monitor Service
After=network.target
Wants=tor.service

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/ping-monitor
Restart=always
RestartSec=5
StandardOutput=append:$INSTALL_DIR/app.log
StandardError=append:$INSTALL_DIR/app.log

[Install]
WantedBy=multi-user.target
EOL

  # Reload systemd
  systemctl daemon-reload
  systemctl enable ping-monitor.service
  
  print_success "ping-monitor set up successfully!"
  print_status "You can now:"
  echo "1. Edit the configuration at $INSTALL_DIR/.env"
  echo "2. Start the service with: systemctl start ping-monitor"
  echo "3. Check logs with: tail -f $INSTALL_DIR/app.log"
else
  print_status "No ping-monitor executable found in current directory."
  print_status "Please build the executable and run this script again."
fi

# Final notes
print_status "Setup complete!"
if [[ $install_tor == "y" || $install_tor == "Y" ]]; then
  print_status "Tor is installed and running. Make sure ENABLE_TOR=true in your .env file."
else
  print_status "If you want to monitor .onion sites, make sure Tor is installed and running."
fi
