.PHONY: build run clean docker docker-run deploy

# Default build target
build:
	go build -o ping-monitor

# Run the application
run: build
	./ping-monitor

# Clean build artifacts
clean:
	rm -f ping-monitor
	rm -rf build
	rm -f ping-monitor.tar.gz

# Build Docker image
docker:
	docker build -t ping-monitor .

# Run in Docker
docker-run: docker
	docker run -d --name ping-monitor \
		-v $(PWD)/.env:/app/.env \
		--restart unless-stopped \
		ping-monitor

# Build deployment package
deploy:
	chmod +x deploy.sh
	./deploy.sh

# Build for production
build-prod:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ping-monitor

# Install as systemd service
install: build-prod
	mkdir -p /opt/ping-monitor
	cp ping-monitor /opt/ping-monitor/
	cp .env /opt/ping-monitor/
	cp ping-monitor.service /etc/systemd/system/
	systemctl daemon-reload
	systemctl enable ping-monitor.service
	systemctl start ping-monitor.service
	@echo "Installed and started ping-monitor service"
