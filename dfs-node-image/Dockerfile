# Start with a lightweight Go image
FROM golang:1.20-alpine

# Set up working directory
WORKDIR /app

# Copy the Go application code
COPY main.go .

# Build the Go binary
RUN go build -o dfs-node main.go

# Create storage directory
RUN mkdir /data

# Expose the HTTP port
EXPOSE 8080

# Set environment variables
ENV NODE_ID="default-node"
ENV STORAGE_CAPACITY="1024"

# Run the node
CMD ["./dfs-node"]
