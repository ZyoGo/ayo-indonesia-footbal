# Build Stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install git and other utilities required for fetching dependencies
RUN apk add --no-cache git tzdata

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o football ./cmd/http/main.go

# Final Stage
FROM alpine:latest

# Add tzdata and ca-certificates
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app/

# Copy the pre-built binary
COPY --from=builder /app/football .

# Copy necessary static assets and configs
COPY --from=builder /app/config ./config
COPY --from=builder /app/keys ./keys

# Expose port (default 4000 unless changed in config)
EXPOSE 4000

# Command to run the executable
CMD ["./football"]
