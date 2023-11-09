# Stage 1: Build the Go binary.
FROM golang:1.21.1 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY src/ src/

# Set the Current Working Directory inside the container
WORKDIR /app/src

# Build the Go app
RUN CGO_ENABLED=0 go build client.go

# Stage 2: Use a distroless base image.
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/src/client .

# Command to run the executable
CMD ["./client"]
