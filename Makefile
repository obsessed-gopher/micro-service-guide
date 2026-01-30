.PHONY: all build run test lint proto docker-build docker-run clean help

# Переменные
APP_NAME := user-service
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go параметры
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Proto
PROTO_DIR := api
PB_DIR := pkg/pb

# Docker
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := $(VERSION)

## help: показать справку по командам
help:
	@echo "Доступные команды:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## all: сборка проекта
all: lint test build

## build: компиляция приложения
build:
	$(GOBUILD) $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/user_service

## run: запуск приложения локально
run:
	$(GOCMD) run ./cmd/user_service -config=config/local.yml

## test: запуск тестов
test:
	$(GOTEST) -v -race -cover ./...

## test-coverage: запуск тестов с отчётом покрытия
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## lint: запуск линтера
lint:
	golangci-lint run ./...

## proto: генерация кода из proto файлов
proto:
	@echo "Генерация proto файлов..."
	@mkdir -p $(PB_DIR)
	protoc --go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/**/*.proto

## deps: загрузка зависимостей
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## docker-build: сборка Docker образа
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f docker/Dockerfile .

## docker-run: запуск в Docker
docker-run:
	docker run --rm -p 50051:50051 $(DOCKER_IMAGE):$(DOCKER_TAG)

## docker-push: отправка образа в registry
docker-push:
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

## clean: очистка артефактов сборки
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

## migrate-up: применение миграций
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

## migrate-down: откат миграций
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

## migrate-create: создание новой миграции
migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name