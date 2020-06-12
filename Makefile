.PHONY: build

start:
	go run cmd/main.go
build:
	env GOOS=linux GOARCH=amd64 go build -o build/rest cmd/main.go