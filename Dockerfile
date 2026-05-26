# Stage 1: Build compilation arena
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency graphs first to cache Docker layer builds efficiently
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the workspace source paths
COPY . .

# Compile a statically linked single production binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/suppliers-binary ./cmd/server/main.go

# Stage 2: Ultra-lightweight secure execution target
FROM alpine:latest  

WORKDIR /root/

# Pull down the compiled asset from the build container stage
COPY --from=builder /app/suppliers-binary .

EXPOSE 8080

# Execute the application
CMD ["./suppliers-binary"]