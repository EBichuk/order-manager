FROM golang:1.23-alpine AS builder

WORKDIR /app

# depen
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN go build -o main ./cmd/api

EXPOSE 8081

CMD ["./cmd/api/main"]