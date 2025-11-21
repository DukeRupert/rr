.PHONY: help build run deploy stop restart logs clean

help:
	@echo "Available commands:"
	@echo "  make build    - Build the Go application locally"
	@echo "  make run      - Run the application locally"
	@echo "  make deploy   - Pull latest code, rebuild Docker image, and start container"
	@echo "  make stop     - Stop and remove running containers"
	@echo "  make restart  - Restart the Docker containers"
	@echo "  make logs     - Follow container logs"
	@echo "  make clean    - Stop containers and remove Docker volumes"

build:
	go build -o main ./cmd/main.go

run:
	go run cmd/main.go

deploy:
	@echo "Pulling latest code..."
	git pull
	@echo "Stopping existing containers..."
	docker-compose down
	@echo "Rebuilding and starting containers..."
	docker-compose up --build -d
	@echo "Deployment complete! Container is running."
	@echo "View logs with: make logs"

stop:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	rm -f main
