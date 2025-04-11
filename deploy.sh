#!/bin/bash

# Build script for go-inactivity-ping
# This script builds the application for Linux and creates a deployment package

# Exit on any error
set -e

# Display build info
echo "Building go-inactivity-ping..."
echo "Go version: $(go version)"

# Set compilation environment for Linux
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# Create build directory
BUILD_DIR="build"
mkdir -p $BUILD_DIR

# Build the application
echo "Compiling for Linux (amd64)..."
go build -ldflags="-s -w" -o $BUILD_DIR/ping-monitor

# Copy necessary files
echo "Creating deployment package..."
cp README.md $BUILD_DIR/
cp .env $BUILD_DIR/  # Copy the actual .env file, not .env.example

# Create empty log file
touch $BUILD_DIR/app.log

# Create a systemd service file
cat > $BUILD_DIR/ping-monitor.service << EOL
[Unit]
Description=URL Monitor Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/ping-monitor
ExecStart=/opt/ping-monitor/ping-monitor
Restart=always
RestartSec=5
StandardOutput=append:/opt/ping-monitor/app.log
StandardError=append:/opt/ping-monitor/app.log

[Install]
WantedBy=multi-user.target
EOL

# Create installation script
cat > $BUILD_DIR/install.sh << EOL
#!/bin/bash
# Installation script for ping-monitor

# Exit on error
set -e

# Create installation directory
INSTALL_DIR="/opt/ping-monitor"
echo "Creating installation directory: \$INSTALL_DIR"
mkdir -p \$INSTALL_DIR

# Copy files
echo "Copying files..."
cp ping-monitor \$INSTALL_DIR/
cp .env \$INSTALL_DIR/
cp README.md \$INSTALL_DIR/
cp app.log \$INSTALL_DIR/
cp ping-monitor.service /etc/systemd/system/

# Set permissions
echo "Setting permissions..."
chmod +x \$INSTALL_DIR/ping-monitor
chmod 644 /etc/systemd/system/ping-monitor.service

# Reload systemd
echo "Reloading systemd..."
systemctl daemon-reload

# Enable service
echo "Enabling service..."
systemctl enable ping-monitor.service

echo "Installation complete!"
echo "Please review your configuration at \$INSTALL_DIR/.env"
echo "Then start the service with: systemctl start ping-monitor"
EOL

# Make install script executable
chmod +x $BUILD_DIR/install.sh

# Create deployment package
echo "Creating tarball..."
tar -czf ping-monitor.tar.gz -C $BUILD_DIR .

echo "Build complete: ping-monitor.tar.gz"
echo "To deploy, transfer this file to your server and extract it"
echo "Then run the install.sh script with sudo permissions"
