# Start from a base Go image
FROM golang:1.19-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the app source code
COPY . .

# Build the Go app
RUN go build -o app .

# Start from a lightweight Debian image
FROM debian:bookworm-slim AS app

# Set the working directory inside the container
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app .

# Run the Go app
CMD ["./app"]
