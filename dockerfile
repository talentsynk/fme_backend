# Use the official Go image as the base
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files (assuming they are in the project root)
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy your source code (replace ./ with your source directory path)
COPY . .

ENV DATABASE_URL=root:subomi7205@tcp(127.0.0.1:3306)/fme_project?charset=utf8mb4&parseTime=True&loc=Local
ENV HASH_SECRET=b21c91b2ab9f6218ce94ee072c43298cb28a

RUN go run cmd/migrations/migration.go
# Build the Go binary (replace main.go with your main file name)
RUN CGO_ENABLED=0 GOOS=linux go build cmd/server/main.go

# Expose the port your application listens on (replace 8080 with your port)
EXPOSE 8080

# Command to run the application
CMD ["./main"]
