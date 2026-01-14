# Используем Go >=1.25
FROM golang:1.23.3-alpine AS builder

# Устанавливаем сертификаты
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Кэш зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Сборка приложения
RUN go build -o myapp ./cmd/monitor

# Минимальный образ для продакшена
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/myapp .
CMD ["./myapp"]
