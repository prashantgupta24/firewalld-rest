.PHONY: build

start:
	go run cmd/main.go
build-linux:
	env GOOS=linux GOARCH=amd64 go build -o build/firewalld-rest cmd/main.go
build-mac:
	go build -o build/firewalld-rest cmd/main.go
clean-db:
	rm /tmp/firewalld-rest-db.tmp
copy: build-linux
	scp build/firewalld-rest root@<server>:/root/rest
