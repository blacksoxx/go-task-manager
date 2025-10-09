FROM golang:1.25.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/api-gateway ./cmd/main.go

FROM golang:1.25.1-alpine
RUN apk add --no-cache ca-certificates

WORKDIR /root/
COPY --from=builder /app/bin/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]