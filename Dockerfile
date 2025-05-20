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

COPY --from=builder /app/ms-authz .
#COPY --from=builder /app/keys ./keys
#COPY --from=builder /app/docs ./docs

RUN apk add --no-cache ca-certificates

EXPOSE 8000

CMD ["./ms-authz"]
