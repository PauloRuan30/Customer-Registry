# Build stage
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api/main.go

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /api /api
EXPOSE 8080
ENTRYPOINT ["/api"]