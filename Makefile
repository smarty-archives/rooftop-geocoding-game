compile:
	mkdir -p static/assets
	GOOS=js GOARCH=wasm go build -o static/main.wasm
	cp assets/* static/assets
serve: compile
	go run http/main.go