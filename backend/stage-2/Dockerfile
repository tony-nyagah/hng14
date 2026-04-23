# ── Builder ───────────────────────────────────────────────────────────────────
FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o server .

# ── Runtime ───────────────────────────────────────────────────────────────────
FROM alpine:latest AS runtime

RUN apk add --no-cache ca-certificates wget && apk upgrade --no-cache

RUN addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -D -s /sbin/nologin appuser

WORKDIR /app

COPY --from=builder /build/server ./server
COPY --from=builder /build/seed_profiles.json ./seed_profiles.json

USER appuser

EXPOSE 8060

HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=3 \
    CMD ["/bin/sh", "-c", "wget -qO- http://localhost:8060/api/profiles?limit=1 > /dev/null 2>&1 || exit 1"]

CMD ["./server"]
