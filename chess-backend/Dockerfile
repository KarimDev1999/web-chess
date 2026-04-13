# Builder stage
FROM golang:1.25-alpine AS builder

# Install migrate
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Generate swagger docs (optional)
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/api/main.go --output docs

# Build the app
RUN go build -o chess-backend ./cmd/api

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy migrate binary and app
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/chess-backend .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

EXPOSE 8080

# Run migrations then start app
CMD ["sh", "-c", "migrate -path /root/migrations -database \"postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable\" up && ./chess-backend"]