FROM --platform=linux/amd64 golang:1.23.0-alpine as builder

WORKDIR /app

COPY . .

RUN go get ./...

RUN RUN go install github.com/a-h/templ/cmd/templ@v0.2.793 && templ generate

RUN go build -tags=jsoniter -o app cmd/main/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/app .

CMD ["./app"]
