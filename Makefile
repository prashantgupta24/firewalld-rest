.PHONY: build

start-local: clean-db
	env=local go run cmd/main.go
start-server:
	go run cmd/main.go
build-linux:
	env GOOS=linux GOARCH=amd64 go build -o build/firewalld-rest cmd/main.go
build-mac:
	go build -o build/firewalld-rest cmd/main.go
copy: build-linux
	scp build/firewalld-rest root@<server>:/root/rest
clean-db:
	rm -f db/*.tmp
test: clean-db
	env=local go test -v -failfast -race ./...
