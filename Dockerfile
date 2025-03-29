# Build Stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for efficient caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Change directory to where the main.go file exists
WORKDIR /app/cmds/server

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

# Final Stage
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /main .

# Ensure it's executable
RUN chmod +x /root/main

EXPOSE 7171

CMD ["/root/main"]
