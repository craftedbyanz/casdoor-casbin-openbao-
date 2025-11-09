.PHONY: run build clean deps test

# Run the server
run:
	go run cmd/server/main.go

# Build the server
build:
	go build -o bin/server cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test ./...

# Start docker services
docker-up:
	docker-compose up -d

# Stop docker services
docker-down:
	docker-compose stop

# View docker logs
docker-logs:
	docker-compose logs -f

# Setup: download deps and start docker
setup: deps docker-up
	@echo "Setup complete!"
	@echo "1. Configure Casdoor application at http://localhost:8000"
	@echo "2. Set CASDOOR_CLIENT_ID and CASDOOR_CLIENT_SECRET in .env"
	@echo "3. Run 'make run' to start the server"

