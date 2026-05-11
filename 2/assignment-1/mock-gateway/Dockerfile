FROM golang:alpine AS builder

WORKDIR /build
COPY ./mock-gateway ./app
WORKDIR /build/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gateway main.go

FROM scratch

WORKDIR /app
COPY --from=builder /build/app/gateway /app/gateway

CMD ["./gateway"]