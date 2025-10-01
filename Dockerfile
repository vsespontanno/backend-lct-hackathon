# builder
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/server

# runtime
FROM gcr.io/distroless/cc-debian11
WORKDIR /app
# Копируем .env файл в контейнер
COPY .env .env
COPY --from=builder /app/bin/app /app/app
EXPOSE 8080
CMD ["/app/app"]