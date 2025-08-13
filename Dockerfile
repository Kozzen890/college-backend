# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o college-backend ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/college-backend .
# COPY --from=builder /app/.env . # (opsional, hanya untuk dev)
EXPOSE 8001
CMD ["./college-backend"]
