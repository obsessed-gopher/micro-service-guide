# syntax=docker/dockerfile:1

# ============================================
# Build stage
# ============================================
FROM golang:1.22-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка приложения
ARG VERSION=dev
ARG BUILD_TIME
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o /app/bin/user-service \
    ./cmd/user_service

# ============================================
# Final stage
# ============================================
FROM alpine:3.19

# Устанавливаем CA сертификаты и timezone данные
RUN apk --no-cache add ca-certificates tzdata

# Создаём непривилегированного пользователя
RUN adduser -D -g '' appuser
USER appuser

WORKDIR /app

# Копируем бинарник из builder stage
COPY --from=builder /app/bin/user-service .

# Копируем конфигурацию
COPY --from=builder /app/config ./config

# Порты
EXPOSE 50051 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Точка входа
ENTRYPOINT ["./user-service"]
CMD ["-config=config/prod.yml"]