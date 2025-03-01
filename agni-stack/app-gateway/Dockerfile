FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go application
RUN go build -o graphql-app ./commander


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/graphql-app .

EXPOSE 8080

# Command to run the executable
CMD ["./graphql-app"]
