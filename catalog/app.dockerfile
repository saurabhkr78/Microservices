# ---------- BUILD STAGE ----------
FROM golang:1.22-alpine AS build

# Disable CGO → static binary
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Install certificates (needed during go mod download)
RUN apk add --no-cache ca-certificates

# Copy go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build your gRPC server
RUN go build -o server ./account/cmd/account


# ---------- RUNTIME STAGE ----------
FROM alpine:3.18

# Install certificates so HTTPS works
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from build stage
COPY --from=build /app/server .

# Security: Run as non-root user
RUN adduser -D appuser
USER appuser

EXPOSE 8080

CMD ["./server"]
