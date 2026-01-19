FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src .
RUN CGO_ENABLED=0 go build -o app ./cmd/server/main.go

FROM alpine:3.8

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
