FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем golang-migrate
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код и собираем приложение
COPY . .
RUN go build -o main .

# Финальный образ
FROM alpine:3.22

WORKDIR /app

# Устанавливаем зависимости для работы с SSL (если нужно)
RUN apk add --no-cache ca-certificates

# Копируем собранное приложение из builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env ./

# Копируем статические файлы если есть
COPY --from=builder /app/storage ./storage

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]