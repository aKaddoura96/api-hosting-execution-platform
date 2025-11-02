# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ shared/

# Copy go mod files for api-gateway
COPY api-gateway/go.mod api-gateway/go.sum api-gateway/
WORKDIR /app/api-gateway

# Download dependencies
RUN go mod download

# Copy api-gateway source code
COPY api-gateway/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/api-gateway/api-gateway .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./api-gateway"]
