COVER_PROFILE=cover.out
COVER_PROFILE_TEMP=cover.tmp.out
COVER_HTML=cover.html

.PHONY: build $(COVER_PROFILE) $(COVER_HTML)

all: coverage vet

coverage: $(COVER_HTML)

$(COVER_HTML): $(COVER_PROFILE) ignoreFiles 
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)

ignoreFiles:
	cat $(COVER_PROFILE_TEMP) | grep -v "middleware.go" | grep -v "route.go" > $(COVER_PROFILE)

$(COVER_PROFILE):
	env=local go test -v -failfast -race -coverprofile=$(COVER_PROFILE_TEMP) ./...

vet:
	go vet ./...
start-local: clean-db #for testing on your local system without firewalld
	env=local go run cmd/main.go
start-server:
	go run cmd/main.go
build-linux: # example: make build-linux DB_PATH=/dir/to/db
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/prashantgupta24/firewalld-rest/db.pathFromEnv=$(DB_PATH)" -o build/firewalld-rest cmd/main.go
local-build:
	go build -ldflags "-X github.com/prashantgupta24/firewalld-rest/db.pathFromEnv=$(DB_PATH)" -o build/firewalld-rest cmd/main.go
copy: build-linux
	scp build/firewalld-rest root@<server>:/root/rest
clean-db:
	rm -f *.db
test:
	env=local go test -v -failfast -race ./...
