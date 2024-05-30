# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./server .

# Stage 2: Create a minimal container
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

ENV GIN_MODE=release

CMD ["/app/server"]
