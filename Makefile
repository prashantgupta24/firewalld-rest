.PHONY: build

start:
	go run cmd/main.go
build-linux:
	env GOOS=linux GOARCH=amd64 go build -o build/rest cmd/main.go
build-local:
	go build -o build/rest cmd/main.go