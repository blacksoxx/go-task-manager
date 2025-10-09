FROM golang:1.25.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/notification-service ./cmd/main.go

FROM golang:1.25.1-alpine
RUN apk add --no-cache ca-certificates

WORKDIR /root/
COPY --from=builder /app/bin/notification-service .

EXPOSE 8083

CMD ["./notification-service"]