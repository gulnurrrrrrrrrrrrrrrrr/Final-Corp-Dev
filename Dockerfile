# Этап 1: Сборка бинарника
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum и скачиваем зависимости (кэшируется)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o quadlingo ./cmd/server

# Этап 2: Финальный минимальный образ
FROM alpine:latest

# Устанавливаем сертификаты (для HTTPS и внешних запросов)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем только бинарник из builder'а
COPY --from=builder /app/quadlingo .

# Копируем статический фронтенд
COPY web/static ./web/static

# Копируем .env (опционально, можно задавать через docker-compose)
COPY .env .

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./quadlingo"]