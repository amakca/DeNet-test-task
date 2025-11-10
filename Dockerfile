## Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build deps (git is often needed by go modules)
RUN apk add --no-cache git

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy only backend sources
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY config ./config
COPY migrations ./migrations

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/cmd/app/app ./cmd/app

## Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

# Layout mirrors repo structure so relative paths in code work:
# - workdir: /app/cmd/app
# - config expected at: ./config/config.yaml
# - migrations expected at: ./../../migrations
WORKDIR /app/cmd/app

# Copy binary and assets
COPY --from=builder /app/cmd/app/app ./app
COPY --from=builder /app/cmd/app/config ./config
COPY --from=builder /app/migrations /app/migrations

ENV HTTP_PORT=8080
EXPOSE 8080

USER app

CMD ["./app"]


