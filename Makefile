bin = todoister

all: build/linux/amd64/$(bin) build/darwin/amd64/$(bin) build/darwin/arm64/$(bin)

build/linux/amd64/$(bin): main.go cmd/*.go util/*.go
	GOOS=linux GOARCH=amd64 go build -o $@

build/darwin/amd64/$(bin): main.go cmd/*.go util/*.go
	GOOS=darwin GOARCH=amd64 go build -o $@

build/darwin/arm64/$(bin): main.go cmd/*.go util/*.go
	GOOS=darwin GOARCH=arm64 go build -o $@

clean:
	rm -rf ./build

.PHONY: clean
