FROM golang:1.24-alpine

RUN apk update && apk upgrade
COPY / .
RUN go build -o main cmd/main.go

CMD ["./main"]

