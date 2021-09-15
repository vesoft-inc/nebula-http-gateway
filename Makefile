.PHONY: build clean fmt

default: build

clean:
	@go mod tidy && rm -rf nebula-httpd

fmt:
	@find . -type f -iname \*.go -exec go fmt {} \;

build: clean fmt
	@go build -o nebula-httpd main.go

run: build
	./nebula-httpd
