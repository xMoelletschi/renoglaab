# Stage 1: Build
FROM golang:1.24@sha256:87916acb3242b6259a26deaa7953bdc6a3a6762a28d340e4f1448e7b5c27c009 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/renoglaab/main.go ./cmd/renoglaab/main.go
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o renoglaab cmd/renoglaab/main.go

# Stage 2: Release
FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715
COPY --from=builder /build/renoglaab /renoglaab
