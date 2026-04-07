# Step 1: Immutable Build Environment
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY . .

# Static compilation: CGO_ENABLED=0 guarantees pure Go binaries that won't require C stdlibs.
# Optimizes for microscopic image footprints and sub-second scale to zero performance.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /api ./cmd/api

# Step 2: Distroless Deployment Output
FROM gcr.io/distroless/static-debian12

WORKDIR /
COPY --from=builder /api /api

EXPOSE 8080
USER nonroot:nonroot

CMD ["/api"]
