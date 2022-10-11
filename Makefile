bin = todo

todo: main.go cmd/*.go
	go build -o $(bin) main.go

clean:
	rm -f $(bin)

.PHONY: clean
