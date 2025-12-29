FROM golang:1.24-bookworm

WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build BOTH binaries
# 1. Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate cmd/migrations/migration.go

# 2. Build the main server
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Expose the internal port
EXPOSE 8000

# Default command (can be overridden by docker-compose)
CMD ["./server"]