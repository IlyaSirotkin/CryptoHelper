FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o run-app ./cmd

FROM alpine
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/run-app .

EXPOSE 8080
CMD ["./run-app"]

