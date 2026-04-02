FROM golang:alpine AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o appointment cmd/appointment/appointment.go

FROM scratch

WORKDIR /app
COPY --from=builder /build/appointment /app/appointment
COPY configs/appointment/config.toml /app/config.toml

CMD ["./appointment"]