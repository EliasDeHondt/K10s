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

COPY App/Backend/go.mod App/Backend/go.sum /app/

RUN go mod tidy

COPY App/Backend/cmd /app/cmd

RUN go build -o /app/main /app/cmd

# Stage 2: Copy the GO app to a smaller image and run it
FROM alpine:3.12 AS production

WORKDIR /app

COPY --from=builder /app/main /app/

EXPOSE 8080

CMD ["./main"]