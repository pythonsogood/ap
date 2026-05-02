FROM golang:alpine AS builder

RUN apk add --no-cache protobuf protobuf-dev git && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && cp /go/bin/protoc-gen-go /usr/bin/

WORKDIR /build
COPY ./appointment-service ./app
COPY ./proto ./proto
RUN protoc --go_out=proto/go --go_opt=paths=source_relative --go-grpc_out=proto/go --go-grpc_opt=paths=source_relative proto/*.proto
WORKDIR /build/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o appointment cmd/appointment/appointment.go

FROM scratch

WORKDIR /app
COPY --from=builder /build/app/appointment /app/appointment
COPY ./appointment-service/configs/appointment/config.toml /app/config.toml
COPY ./appointment-service/migrations /app/migrations

CMD ["./appointment"]