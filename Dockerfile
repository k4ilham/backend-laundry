FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application pointing to the correct main file
RUN go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /app/main .

# Expose port (matched with SERVER_PORT in railway env)
EXPOSE 8090

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
    CMD curl -f http://localhost:8090/health || exit 1

CMD ["./main"]
