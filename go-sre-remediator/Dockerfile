# --- Stage 1: Build the Go Binary ---
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker caching for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with optimizations for containers (disable CGO, target Linux)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o remediator main.go

# --- Stage 2: Final Lightweight Runtime ---
FROM alpine:3.19

WORKDIR /app

# Install bash since Alpine uses bash by default, and our engine runs /bin/bash
RUN apk add --no-cache bash ca-certificates

# 1. Create your unprivileged system group and user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup -u 100001

# 2. Copy the compiled binary from the builder stage
COPY --from=builder /app/remediator .

# 3. Copy our configuration and remediation scripts
COPY config.yaml .
COPY remediate_disk.sh .
COPY remediate_service.sh .

# 4. Ensure scripts have execution permissions inside the container
RUN chmod +x *.sh && chown -R appuser:appgroup /app

# 5. Explicitly switch security contexts away from root
USER 10001

# Expose the ports for webhooks and Prometheus scraping
EXPOSE 8080

# Run the engine
CMD ["./remediator"]