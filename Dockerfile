# Stage 1: Build
FROM golang:1.25@sha256:698183780de28062f4ef46f82a79ec0ae69d2d22f7b160cf69f71ea8d98bf25d AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/renoglaab/main.go ./cmd/renoglaab/main.go
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o renoglaab cmd/renoglaab/main.go

# Stage 2: Release
FROM alpine:3.23@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62
COPY --from=builder /build/renoglaab /renoglaab
