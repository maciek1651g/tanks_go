# Use the official Go image as the base image
FROM golang:1.20

# Set the working directory in the container
WORKDIR ./src

# Copy the application files into the working directory
COPY ./ ./

# Build the application
RUN go build -o main ./src

# Expose port 8080
EXPOSE 8080

# Define the entry point for the container. The entry point(command) should start the application.
CMD ["./main"]