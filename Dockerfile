FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


WORKDIR  /app/src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /out/assignment-service



FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN adduser -D -H -u 10001 appuser

COPY --from=builder /out/assignment-service /usr/bin/assignment-service

RUN printf '%s\n' \
  '#!/bin/sh' \
  'umask 0027' \
  'exec /usr/bin/assignment-service "$@"' > /usr/local/bin/entrypoint && \
  chmod +x /usr/local/bin/entrypoint

WORKDIR /app
USER appuser

ENTRYPOINT ["/usr/local/bin/entrypoint"]