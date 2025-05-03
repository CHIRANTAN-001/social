# Builder stage
FROM golang:1.23.8-alpine3.20 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o ./bin/main ./cmd/api/

# Runtime stage
FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/bin/main ./bin/main
RUN chmod +x ./bin/main

EXPOSE 8080

CMD ["./bin/main"]