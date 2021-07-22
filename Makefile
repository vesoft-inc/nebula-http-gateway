.PHONY: build clean fmt test

default: build

clean:
	@go mod tidy && rm -rf nebula-httpd

fmt:
	@find . -type f -iname \*.go -exec go fmt {} \;

test: clean fmt
	@go test -v ./...

build: clean fmt
	@go build -o nebula-httpd main.go
