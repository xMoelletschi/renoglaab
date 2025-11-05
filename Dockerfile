# Stage 1: Build
FROM golang:1.24@sha256:62c7ec2059c165fee0750d3ee51429d6c167749f70dd13b7dfea2c85dfd941de AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/renoglaab/main.go ./cmd/renoglaab/main.go
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o renoglaab cmd/renoglaab/main.go

# Stage 2: Release
FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715
COPY --from=builder /build/renoglaab /renoglaab
