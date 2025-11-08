FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY common/ ./common/
COPY credit-exchange-service/ ./credit-exchange-service/

WORKDIR /app/credit-exchange-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/credit-exchange-service/main .

EXPOSE 8086

CMD ["./main"]
