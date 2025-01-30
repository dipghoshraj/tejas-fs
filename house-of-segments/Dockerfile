

FROM golang:1.22.2-alpine AS builder


WORKDIR /app
COPY go.mod go.sum ./


RUN go mod tidy

COPY . .

RUN go build -o grpc-server ./unit-commander


FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/grpc-server /usr/local/bin/grpc-server

# Step 10: Expose the port that your gRPC server will be listening on
EXPOSE 50051

# Step 11: Set the entrypoint to the application
CMD ["grpc-server"]
