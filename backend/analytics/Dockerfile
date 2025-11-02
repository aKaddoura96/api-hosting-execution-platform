# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ shared/

# Copy go mod files for analytics
COPY analytics/go.mod analytics/go.sum analytics/
WORKDIR /app/analytics

# Download dependencies
RUN go mod download

# Copy analytics source code
COPY analytics/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o analytics main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/analytics/analytics .

# Expose port
EXPOSE 8082

# Run the application
CMD ["./analytics"]
