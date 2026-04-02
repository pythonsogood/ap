FROM golang:alpine AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o doctor cmd/doctor/doctor.go

FROM scratch

WORKDIR /app
COPY --from=builder /build/doctor /app/doctor
COPY configs/doctor/config.toml /app/config.toml

CMD ["./doctor"]