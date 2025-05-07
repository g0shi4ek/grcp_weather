# Этап сборки
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Переходим в директорию сервера и собираем
WORKDIR /app/cmd/startClient
RUN CGO_ENABLED=0 GOOS=linux go build -o startClient .

# Этап запуска
FROM alpine:latest
WORKDIR /app

# Копируем собранный бинарник и зависимости
COPY --from=builder /app/cmd/startClient/startClient .

CMD ["./startClient"]