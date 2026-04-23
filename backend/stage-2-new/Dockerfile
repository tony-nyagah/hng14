# Stage 1: Build the binary
FROM golang:1.26-alpine AS builder

# Install build dependencies for sqlite3 (requires CGO)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the app with CGO enabled for sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Stage 2: Final lightweight image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
# Copy seed data so the app can seed on startup
COPY --from=builder /app/seed.json .

# Expose the port the app runs on
EXPOSE 8070

# Run the application
CMD ["./main"]
