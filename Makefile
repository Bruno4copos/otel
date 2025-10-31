# Nome do projeto
PROJECT_NAME=servico-cep-temperatura

# Diretórios dos serviços
SERVICE_A_DIR=servico-a
SERVICE_B_DIR=servico-b

# Imagens Docker
SERVICE_A_IMAGE=$(PROJECT_NAME)-a
SERVICE_B_IMAGE=$(PROJECT_NAME)-b

# OTEL Collector
OTEL_IMAGE=otel/opentelemetry-collector:latest

# Comandos padrão
GO_BUILD=go build -o bin/
GO_RUN=go run
DOCKER_COMPOSE=docker-compose -f docker-compose.yaml

.PHONY: all build run clean up down logs fmt lint test

## ========================
## BUILD & DEV
## ========================

build:
	@echo "🔧 Building Service A..."
	cd $(SERVICE_A_DIR) && $(GO_BUILD)
	@echo "🔧 Building Service B..."
	cd $(SERVICE_B_DIR) && $(GO_BUILD)
	@echo "✅ Build complete!"

run-a:
	@echo "🚀 Running Service A locally..."
	cd $(SERVICE_A_DIR) && $(GO_RUN) .

run-b:
	@echo "🚀 Running Service B locally..."
	cd $(SERVICE_B_DIR) && $(GO_RUN) .

fmt:
	@echo "🧹 Formatting code..."
	go fmt ./...

lint:
	@echo "🔍 Linting project..."
	golangci-lint run ./...

test:
	@echo "🧪 Running unit tests..."
	go test ./... -v

clean:
	@echo "🧼 Cleaning build artifacts..."
	rm -rf $(SERVICE_A_DIR)/bin $(SERVICE_B_DIR)/bin

## ========================
## DOCKER / DEPLOY
## ========================

up:
	@echo "🛠️  Starting full environment (Services + OTEL + Zipkin)..."
	$(DOCKER_COMPOSE) up --build -d
	@echo "✅ All services are up!"

down:
	@echo "🛑 Stopping all containers..."
	$(DOCKER_COMPOSE) down

logs:
	@echo "📜 Showing logs (Ctrl+C to exit)..."
	$(DOCKER_COMPOSE) logs -f

## ========================
## HELP
## ========================

help:
	@echo "Available commands:"
	@echo "  make build    - Build both Go services"
	@echo "  make run-a    - Run Service A locally"
	@echo "  make run-b    - Run Service B locally"
	@echo "  make up       - Start all containers (services + OTEL + Zipkin)"
	@echo "  make down     - Stop and remove containers"
	@echo "  make logs     - Tail logs from all services"
	@echo "  make clean    - Clean binaries"
