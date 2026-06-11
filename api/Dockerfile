FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /api .

# PocketBase stores its data (SQLite DB, uploads, etc.) in this directory
VOLUME /app/pb_data

EXPOSE 8090
CMD ["./api", "serve", "--http=0.0.0.0:8090", "--dir=/app/pb_data"]
