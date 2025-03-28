# Use an updated base image with GLIBC 2.33+
FROM ubuntu:22.04 AS builder

WORKDIR /app

# Install dependencies
RUN apt update && apt install -y build-essential curl

# Install Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Copy project files
COPY . .

# Build the Rust binary
RUN cargo build --release

# Use the same updated base image for runtime
FROM ubuntu:22.04

WORKDIR /app

# Install required runtime libraries
RUN apt update && apt install -y libc6 && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary
COPY --from=builder /app/target/release/clusterapp /app/clusterapp

CMD ["./clusterapp"]
