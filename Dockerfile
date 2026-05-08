# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o phone .

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary and static files from builder
COPY --from=builder /app/phone .
COPY --from=builder /app/index.html .

# Expose the web server port
EXPOSE 8080

# Run the application
CMD ["./phone"]
