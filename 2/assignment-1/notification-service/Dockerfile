FROM golang:alpine AS builder

WORKDIR /build
COPY ./notification-service ./app
WORKDIR /build/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o notification cmd/notification/notification.go

FROM scratch

WORKDIR /app
COPY --from=builder /build/app/notification /app/notification
COPY ./notification-service/configs/notification/config.toml /app/config.toml

CMD ["./notification"]