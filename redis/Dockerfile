# Build stage
FROM golang:1.25-alpine AS builder

RUN adduser -D -u 1000 appuser
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o redis102

# Run stage
FROM scratch

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

USER appuser
COPY --from=builder /app/redis102 /redis102

# Expose port (default in config is 8080)
EXPOSE 8080

# Command to run the executable
ENTRYPOINT [ "/redis102" ]
