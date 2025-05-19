# ==== STAGE 1: Build binary ====
FROM golang:1.21-alpine AS builder

# Install CA certs (for HTTP requests, optional)
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux go build -o ms-authz ./cmd/main.go


# ==== STAGE 2: Create minimal runtime container ====
FROM alpine:latest

# Create working dir
WORKDIR /app

# Copy built binary and keys
COPY --from=builder /app/ms-authz .
COPY --from=builder /app/keys ./keys

# Add CA certs (in case public key loading needs HTTPS)
RUN apk add --no-cache ca-certificates

# Expose app port
EXPOSE 8000

# Run
CMD ["./ms-authz"]
