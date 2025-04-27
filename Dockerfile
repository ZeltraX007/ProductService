# Start from the official Golang base image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o productservice

# Expose the port your service will run on
EXPOSE 8000

# Command to run the executable
CMD ["./productservice"]