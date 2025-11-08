# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates wget

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service ./notification-service/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/notification-service .

EXPOSE 8087

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8087/health/live || exit 1

CMD ["./notification-service"]
