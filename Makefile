.PHONY: build run test clean

build:
	go build -o bin/diamond ./cmd/diamond

run: build
	./bin/diamond

test:
	go test ./...

clean:
	rm -rf bin
