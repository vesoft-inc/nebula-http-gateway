.PHONY: build clean fmt

default: build

clean:
	@go mod tidy && rm -rf nebula-httpd

build: clean fmt
	@go build -o nebula-httpd main.go

fmt:
	@find . -type f -iname \*.go -exec go fmt {} \;
