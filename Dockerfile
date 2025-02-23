# Stage 1: Build
FROM golang:1.24@sha256:2b1cbf278ce05a2a310a3d695ebb176420117a8cfcfcc4e5e68a1bef5f6354da AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/renoglaab/main.go ./cmd/renoglaab/main.go
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o renoglaab cmd/renoglaab/main.go

# Stage 2: Release
FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c
COPY --from=builder /build/renoglaab /renoglaab
