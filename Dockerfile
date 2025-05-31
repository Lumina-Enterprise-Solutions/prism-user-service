FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o prism-user-service ./cmd/server/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/prism-user-service .
EXPOSE 8080
CMD ["./prism-user-service"]
