# -----------------------------
# Build Stage
# -----------------------------
FROM golang:1.24.0-alpine3.20 AS build

# Install required tools
RUN apk --no-cache add gcc g++ make ca-certificates

# Set working directory
WORKDIR /github.com/saurabh/Microservices

# Copy module files
COPY go.mod go.sum ./

# Copy vendor + services
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order

# Build the binary
RUN go build -mod=vendor -o /go/bin/app ./order/cmd/order

# -----------------------------
# Runtime Stage
# -----------------------------
FROM alpine:3.20

WORKDIR /usr/bin

# Copy only the built binary
COPY --from=build /go/bin/app .

EXPOSE 8080
CMD ["./app"]
