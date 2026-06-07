.PHONY: build test clean

build:
	go build -o bin/rook .

test:
	go test ./... -v -count=1

clean:
	rm -rf bin/
