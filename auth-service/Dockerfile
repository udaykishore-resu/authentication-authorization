FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# First copy module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build
# RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go
RUN go build -o auth-service ./cmd/main.go
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/auth-service .
COPY --from=builder /app/config/config.yaml ./config/
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/private.pem .
COPY --from=builder /app/public.pem .

EXPOSE 8080
CMD ["./auth-service"]