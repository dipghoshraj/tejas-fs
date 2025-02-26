

FROM golang:1.22.2-alpine AS builder


WORKDIR /app
COPY go.mod go.sum ./


RUN go mod tidy

COPY . .

RUN go build -o grpc-server ./unit-commander


FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/grpc-server /usr/local/bin/grpc-server

EXPOSE 50051

CMD ["grpc-server"]
