FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@v0.2.793 && templ generate

RUN go get ./...

RUN go build -tags=jsoniter -o app cmd/main/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/app .

CMD ["./app"]
