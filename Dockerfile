FROM golang:1.26.4-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o edgetpu-exporter .

FROM alpine:3.24

RUN addgroup -S exporter && adduser -S -u 1000 -G exporter exporter

COPY --from=builder /app/edgetpu-exporter /app/edgetpu-exporter

USER exporter

ENTRYPOINT ["/app/edgetpu-exporter"]
