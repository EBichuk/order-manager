FROM golang:1.23-alpine AS builder

WORKDIR /app/producer

# depen
COPY /producer/go.mod /producer/go.sum ./
RUN go mod download

# build
COPY producer/ ./
RUN go build -o app ./main.go

CMD ["./app"]