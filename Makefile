bin = todo

todo: main.go cmd/root.go
	go build -o $(bin) main.go

clean:
	rm -f $(bin)

.PHONY: clean
