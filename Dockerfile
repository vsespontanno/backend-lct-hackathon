# Этап сборки
FROM golang:1.24.0-alpine AS builder

WORKDIR /app

# Скопировать go.mod и go.sum, подтянуть зависимости
COPY go.mod go.sum ./
RUN go mod download

# Скопировать исходники
COPY . .

# Собрать бинарник
RUN go build -o main ./cmd/server

# Этап запуска
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/main .

# Если нужно, подтянем .env внутрь контейнера
COPY .env .env

EXPOSE 8080
CMD ["./main"]