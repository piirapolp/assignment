FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache --virtual .build-deps ca-certificates

WORKDIR /app
COPY . /app

WORKDIR  /app/src
RUN go build -o ../build/assignment-service



FROM alpine:3.20

RUN umask 027 && echo "umask 0027" >> /etc/profile

COPY --from=builder /var/cache/apk /var/cache/apk

RUN apk add --no-cache bash tzdata ca-certificates && rm -rf /var/cache/apk

COPY --from=builder /app/build/assignment-service /usr/bin/assignment-service

WORKDIR /app
ENTRYPOINT [ "assignment-service" ]