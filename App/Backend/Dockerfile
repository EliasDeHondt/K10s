#################################
# @author K10s Open Source Team #
# @since 01/01/2025             #
#################################
# Backend build
# Stage 1: Build the GO app
FROM golang:1.20 AS builder

LABEL maintainer="K10s Open Source Team"
LABEL version="1.0"
LABEL description="This is the backend image for the app."
LABEL org.opencontainers.image.description="This is the backend image for the app."

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Stage 2: Copy the GO app to a smaller image and run it
FROM alpine:3.12 AS production

WORKDIR /app

COPY --from=builder /app/main /app/

EXPOSE 8080

CMD ["./main"]