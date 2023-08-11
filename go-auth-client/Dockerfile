# Use the official Go image as the base image
FROM golang:1.20.6

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the source code and template files into the container
COPY . .

# Build the Go app
RUN go build -o main

# Set the command to run the Go app
CMD ["./main"]
