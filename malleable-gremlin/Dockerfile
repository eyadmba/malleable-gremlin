# ---- Builder Stage ----
# Use an official Go image as the builder
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# Assuming your main package is in cmd/server
# Statically link the binary and strip debug information
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/malleable-gremlin ./server/

# ---- Final Stage ----
# Use a minimal Alpine image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/malleable-gremlin /app/malleable-gremlin

# Environment variable for the port, defaulting to 8080
ENV PORT=8080

# Expose the default port (for documentation and potential linking)
# The actual port mapping is done via `docker run -p`
EXPOSE 8080

# Command to run the application
# The application should listen on the address defined by the PORT env var (e.g., ":8080")
ENTRYPOINT ["/app/malleable-gremlin"]