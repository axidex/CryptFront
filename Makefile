run: template run-server

run-server:
	go run cmd/main/main.go

tidy:
	go mod tidy
	go fmt ./...


template:
	templ generate