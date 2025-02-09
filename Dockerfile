FROM node:18-alpine AS css-builder


WORKDIR /app

# COPY ./static/css ./static/css
# COPY ./tailwind.config.js ./
# COPY ./views ./views

COPY . .

RUN npm install -D tailwindcss@3
RUN npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css

FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@v0.2.793 && templ generate

RUN go get ./...

COPY --from=css-builder /app/static/css/output.css ./static/css/

RUN go build -tags=jsoniter -o app cmd/main/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=builder /app/app .
COPY --from=css-builder /app/static/css/output.css ./static/css/

CMD ["./app"]
