data:
	cat .satisfy | satisfy
compile: data
	mkdir -p static/assets
	GOOS=js GOARCH=wasm go build -o static/main.wasm
	cp ~/.cache/satisfy/platformer/* static/assets
	cp files/* static/
serve: compile
	go run http/main.go