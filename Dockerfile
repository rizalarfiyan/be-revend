# Build stage
FROM golang:1.21.1-alpine AS builder
WORKDIR /app
COPY . .

RUN apk add --no-cache curl && \
    rm -rf /var/cache/apk/*

# Download goose binary
RUN curl -SL https://github.com/pressly/goose/releases/download/v3.15.0/goose_linux_x86_64 -o goose && \
    chmod +x goose

# Download swag binary and generate swagger docs
RUN curl -SL https://github.com/swaggo/swag/releases/download/v1.16.2/swag_1.16.2_Linux_x86_64.tar.gz -o swag.tar.gz && \
    mkdir -p ./swag && \
    tar -xzf swag.tar.gz --directory ./swag && \
    rm swag.tar.gz && \
    ./swag/swag init && \
    rm -rf ./swag

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
COPY --from=builder /app/sql /app/sql
COPY --from=builder /app/migrate.sh /app/
COPY --from=builder /app/goose /app/

CMD ["/app/be-revend"]
