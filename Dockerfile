FROM golang:1.25.1-alpine AS builder-migrator

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /build/app ./cmd/api/

FROM alpine:3.22 AS app

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /app/migrations ./migrations/

EXPOSE 8081

CMD ["./app"]

FROM alpine:3.22 AS goose

COPY --from=builder-migrator /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations /migrations

WORKDIR /migrations

ENTRYPOINT ["goose"]

CMD ["up"]