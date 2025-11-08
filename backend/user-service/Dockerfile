FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY common/ ./common/
COPY user-service/ ./user-service/

WORKDIR /app/user-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/user-service/main .

EXPOSE 8080

CMD ["./main"]
