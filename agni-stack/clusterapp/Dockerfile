FROM rust:latest as builder
WORKDIR /app
COPY . .
RUN cargo build --release
RUN cargo install --path .


FROM debian:bullseye-slim
WORKDIR /app

COPY --from=builder /app/target/release/clusterapp /app/clusterapp
RUN apt-get update && apt-get install -y libssl-dev libpq-dev ca-certificates  && rm -rf /var/lib/apt/lists/*
CMD ["./clusterapp"]