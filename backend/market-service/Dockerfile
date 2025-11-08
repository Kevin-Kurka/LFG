FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY common/ ./common/
COPY market-service/ ./market-service/

WORKDIR /app/market-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/market-service/main .

EXPOSE 8082

CMD ["./main"]
