# Build stage
FROM golang:1.20-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o probing .

# Final stage
FROM scratch

# Copy the statically compiled binary from the builder stage
COPY --from=builder /app/probing /probing

# Set the entrypoint
ENTRYPOINT ["/probing"]
