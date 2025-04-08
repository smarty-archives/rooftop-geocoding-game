data:
	mkdir -p static
	rm -r static/
	cat .satisfy | satisfy
	mkdir -p static/assets
	cp -r ~/.cache/satisfy/platformer/* static/assets
	cp files/* static/
compile: data
	mkdir -p assets
	rm -r assets
	cp -r static/assets .
	GOOS=js GOARCH=wasm go build -o static/main.wasm
serve: compile
	go run http/main.go
run: compile
	go run *.go