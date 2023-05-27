FROM golang:1.19-alpine as builder

WORKDIR /app

COPY . /app
RUN go mod tidy
RUN go build -o /app main.go

FROM alpine:3
WORKDIR /app

RUN mkdir -p /var/local_storage/file_server/

COPY config.yaml .
COPY private.pem .
COPY --from=builder /app/main /app
