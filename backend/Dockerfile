# Development stage with Air for hot reloading
FROM golang:1.22-alpine AS dev

WORKDIR /app

# Install dependencies for development
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Install Air for hot reloading
RUN go install github.com/cosmtrek/air@v1.49.0

# Copy source code
COPY . .

# Create tmp directory for Air
RUN mkdir -p tmp

# Expose port
EXPOSE 8080

# Command to run Air for hot reloading
CMD ["air", "-c", ".air.toml"]

# Production builder stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Production stage
FROM alpine:latest AS prod

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
