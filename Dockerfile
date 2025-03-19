# Start from a base Go image
FROM golang:1.23-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the app source code
COPY main.go .

# Set environment variables to build a statically linked binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build the Go app
RUN go build -o app .

# Start from a scratch image
FROM scratch AS prod

WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/app .
COPY calendar_template.png .
COPY Roboto-Bold.ttf .

ENV GIN_MODE=release

CMD ["./app"]