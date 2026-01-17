
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY ./src .

RUN go mod download

RUN CGO_ENABLED=0 go build -o app ./cmd/server/main.go

FROM alpine:3.8

WORKDIR /app
COPY --from=builder app .

EXPOSE 8080
CMD ["./app"]