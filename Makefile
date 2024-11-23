BIN = todoister
VERSION = 0.2.0

TAG = $(shell git describe --tags --always --abbrev=0)
LDFLAGS= -ldflags="-X 'github.com/layfellow/todoister/cmd.Version=$(VERSION)'"

build:
	go build $(LDFLAGS) -o build/$(BIN)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/$(BIN)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/$(BIN)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o build/$(BIN)-darwin-arm64

lint:
	golangci-lint run

dependencies:
	go get -u
	go mod tidy

releases:
	gh release create $(TAG) ./build/$(BIN)-linux-amd64 ./build/$(BIN)-darwin-amd64 ./build/$(BIN)-darwin-arm64

install:
	go env -w GOBIN=$$HOME/bin
	go install

clean:
	rm -rf build

.PHONY: build lint dependencies releases install clean
