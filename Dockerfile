# Use the official Go image to build the app
FROM golang:1.23.3 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go module files into the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire app into the container
COPY . .

# Set the directory where your `main.go` file resides
WORKDIR /app

# Build the Go binary
RUN go build -o /app/app

# Stage 2: Create the final image for running the app
FROM golang:1.23.3 

# Set the environment variables
ENV DATABASE_PATH=data/databases/production

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/app .

# Copy the static and data directories
COPY --from=builder /app/static /root/static
COPY --from=builder /app/data /root/data

# Expose the port the app will run on
EXPOSE 8080

# Run the Go app
CMD ["/root/app"]