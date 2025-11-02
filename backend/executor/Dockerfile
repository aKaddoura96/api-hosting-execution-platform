# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ shared/

# Copy go mod files for executor
COPY executor/go.mod executor/
WORKDIR /app/executor

# Download dependencies
RUN go mod download

# Copy executor source code
COPY executor/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o executor -mod=mod main_simple.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/executor/executor .

# Expose port
EXPOSE 8081

# Run the application
CMD ["./executor"]
