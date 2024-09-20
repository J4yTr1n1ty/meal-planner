# Start from the official Go image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Install gcc and stdlib for CGO
RUN apk add --no-cache gcc libc-dev

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /server cmd/server/main.go

# Start a new stage from scratch
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /server .

# Copy the static files
COPY --from=builder /app/static ./static

# Command to run the executable
CMD ["./server"]

