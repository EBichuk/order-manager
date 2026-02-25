FROM golang:1.25.1-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .

RUN go build -o producer main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /build/producer .

CMD ["./producer"]