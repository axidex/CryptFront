run: template run-server

run-server:
	go run cmd/main/main.go

tidy:
	go mod tidy
	go fmt ./...


template:
	templ generate

run-templ:
	templ generate --watch --proxy="http://localhost:5050" --cmd="go run cmd/main/main.go"


