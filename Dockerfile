
FROM golang:1.25-alpine As builder
WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod .
COPY go.sum .

RUN go mod download

