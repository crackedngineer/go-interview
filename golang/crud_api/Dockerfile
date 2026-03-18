# Build stage
FROM golang:1.25-bookworm AS builder

# Install gcc for CGO
RUN apt-get update && apt-get install -y gcc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Enable CGO
ENV CGO_ENABLED=1

RUN go build -o app cmd/student_api/main.go

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

# Optional: copy config if needed
COPY config ./config

CMD ["./app"]