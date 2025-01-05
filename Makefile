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

styles:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css

# Run templ generation in watch mode
templ:
	templ generate --watch --proxy="http://localhost:8090" --open-browser=false -v

# Run air for Go hot reload
server:
	air \
    --build.cmd "go build -o tmp/main.exe cmd/main/main.go" \
    --build.bin "tmp\main.exe" \
    --build.delay "100" \
    --build.exclude_dir ".idea,node_modules,tmp" \
    --build.include_ext "go" \
    --build.stop_on_error "true" \
    --misc.clean_on_exit true \
	--build.kill_delay "3s" \
	--root .

# Watch Tailwind CSS changes
tailwind:
	tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

# Start development server with all watchers
dev:
	make -j3 templ server tailwind