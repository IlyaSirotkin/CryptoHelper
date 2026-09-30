FROM golang:1.27-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /src/main \
    cmd/main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

RUN addgroup -g 1000 nonroot \
    && adduser -D -u 1000 -G nonroot -s /sbin/nologin nonroot

WORKDIR /app

RUN mkdir -p /app/.log \
    && chown -R nonroot:nonroot /app

COPY --from=builder --chown=nonroot:nonroot /src/main /app/main

USER nonroot:nonroot

ENTRYPOINT ["./main"]
