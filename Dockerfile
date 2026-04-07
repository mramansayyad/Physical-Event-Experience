# Use the latest possible Go 1.25 image
FROM golang:1.25-bookworm AS builder

# Install build essentials
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Force the Go toolchain to auto-upgrade to 1.25 if needed
ENV GOTOOLCHAIN=go1.25.0

# Copy everything
COPY . .

# Clean dependencies - this is where the 1.25 requirement was failing
RUN rm -f go.sum
RUN go mod tidy

# Build the binary
RUN go build -v -o /stadium-backend ./cmd/api/*.go

# STEP 2: Final Image (Standard Debian Slim for maximum compatibility)
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /
COPY --from=builder /stadium-backend /stadium-backend

ENV PORT=8080
EXPOSE 8080

CMD ["/stadium-backend"]