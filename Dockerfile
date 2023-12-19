# Build stage
FROM golang:1.21.1-alpine AS builder
WORKDIR /app
COPY . .

# Build the Go application
RUN go mod download
RUN go mod verify
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o be-revend

# Final stage
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache tzdata && \
    rm -rf /var/cache/apk/*

# Copy the binary and other necessary files
COPY --from=builder /app/be-revend /app/

CMD ["/app/be-revend"]
