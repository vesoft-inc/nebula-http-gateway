
.PHONY: build clean

default: build

clean:
	go mod tidy && rm -rf nebula-httpd

build: clean
	go build -o nebula-httpd main.go
