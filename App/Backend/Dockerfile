#################################
# @author K10s Open Source Team #
# @since 01/01/2025             #
#################################
# Backend build
# Step 1: Build the Go binary
FROM golang:1.23-alpine as builder

LABEL maintainer "K10s Open Source Team"
LABEL version "1.0"
LABEL description "This is the backend image for the app."
LABEL org.opencontainers.image.description "This is the backend image for the app."

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./cmd

# Step 2: Create a smaller image for running the app
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]