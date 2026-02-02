# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod (and go.sum if exists)
COPY go.mod ./
# COPY go.sum ./ 

# Download dependencies
# RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Assumes main.go is in root for answer-service initially
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run Stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
