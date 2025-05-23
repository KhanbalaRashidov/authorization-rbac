# ==== STAGE 1: Build binary ====
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ms-authz ./cmd/main.go


# ==== STAGE 2: Create minimal runtime container ====
FROM alpine:latest

WORKDIR /app

ARG DB_HOST
ARG DB_NAME
ARG DB_USERNAME
ARG DB_PASSWORD
ARG RABBITMQ_USER
ARG RABBITMQ_PASSWORD
ARG RABBITMQ_HOST

ENV DB_DSN="host=$DB_HOST user=$DB_USERNAME password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable"
ENV RABBITMQ_URL=amqp://$RABBITMQ_USER:$RABBITMQ_PASSWORD@$RABBITMQ_HOST:5672/
ENV PUBLIC_KEY_DIR=/app/keys
ENV PORT=8080

COPY --from=builder /app ./
#COPY --from=builder /app/keys ./keys
#COPY --from=builder /app/docs ./docs

RUN apk add --no-cache ca-certificates

EXPOSE 8080

CMD ["./ms-authz"]
