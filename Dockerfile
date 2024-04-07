# Use a Golang base image
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project directory into the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Expose the port on which the Go application will listen
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
