bin = todo

todo: main.go cmd/root.go cmd/version.go
	go build -o $(bin) main.go

clean:
	rm -f $(bin)

.PHONY: clean
