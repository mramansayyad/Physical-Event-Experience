# Build Stage
FROM golang:1.25.9-alpine AS builder

# Enable CGO_ENABLED=0 for static binaries
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOTOOLCHAIN=go1.25.9

WORKDIR /app

# Optimize caching of go.mod / go.sum
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Build the binary statically, stripping symbols and injecting LDFlags for version dynamically
RUN go build -ldflags="-s -w -X 'main.Version=v1.1.0-SECURED'" -o /stadium-backend ./cmd/api/*.go

# Final Stage (Distroless for Zero-Trust)
# Uses a minimal base with NO shell or package managers.
FROM gcr.io/distroless/static-debian12:nonroot

# Use the non-root user predefined in the distroless image
USER nonroot:nonroot

COPY --from=builder /stadium-backend /stadium-backend

ENV PORT=8080
EXPOSE 8080

CMD ["/stadium-backend"]