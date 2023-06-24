# Use the official Golang image as the base image
FROM golang:1.19-alpine

# Create a working directory
WORKDIR /app

# Copy the source code into the working directory
COPY . .

RUN go mod download

# Build the Golang application
RUN go build -o main .

# Set the default command to run when the container starts
CMD ["./main"]
