# Use an official Golang runtime as the base image
FROM golang:1.21.6

# Set the working directory inside the container
WORKDIR /app

# Copy mods and sum files
COPY go.mod .
COPY go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Expose the port your application will listen on
EXPOSE 8080

# Define the command to run your application
CMD ["./main"]
