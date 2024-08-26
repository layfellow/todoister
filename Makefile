bin = todoister

all: build/linux/amd64/$(bin) build/darwin/amd64/$(bin) build/darwin/arm64/$(bin) build/windows/amd64/$(bin).exe

build/linux/amd64/$(bin): main.go cmd/*.go
	GOOS=linux GOARCH=amd64 go build -o $@ main.go

build/darwin/amd64/$(bin): main.go cmd/*.go
	GOOS=darwin GOARCH=amd64 go build -o $@ main.go

build/darwin/arm64/$(bin): main.go cmd/*.go
	GOOS=darwin GOARCH=arm64 go build -o $@ main.go

build/windows/amd64/$(bin).exe: main.go cmd/*.go
	GOOS=windows GOARCH=amd64 go build -o $@ main.go

clean:
	rm -rf ./build

.PHONY: clean
