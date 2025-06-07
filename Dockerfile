# ---- Build Stage ----
FROM golang:1.23-alpine AS builder

# Install git (required for go get) and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN go build -o firelink main.go

# ---- Final Stage ----
FROM alpine:latest

# Install ca-certificates and pg_dump (for database backup)
RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/firelink .

# Expose the port your app runs on
EXPOSE 8080

# Set environment variables (optional, can also be set in docker-compose)
# ENV POSTGRES_HOST=postgres
# ENV POSTGRES_USER=yourusername
# ENV POSTGRES_PASSWORD=yourpassword
# ENV POSTGRES_DB=firelink

# Run the binary
CMD ["./firelink"]
