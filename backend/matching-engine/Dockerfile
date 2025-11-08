FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY common/ ./common/
COPY matching-engine/ ./matching-engine/

WORKDIR /app/matching-engine

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/matching-engine/main .

EXPOSE 8084

CMD ["./main"]
