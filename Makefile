data:
	mkdir -p static
	rm -r static/
	cat .satisfy | satisfy
	mkdir -p static/assets
	cp -r ~/.cache/satisfy/platformer/* static/assets
	cp files/* static/
compile: data
	GOOS=js GOARCH=wasm go build -o static/main.wasm
serve: compile
	go run http/main.go
run: data
	mkdir -p assets
	rm -r assets
	cp -r static/assets .
	go run main_local.go