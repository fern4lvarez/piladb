.PHONY: default all test testv race vet build clean

default: vet test

all: vet test build

test:
	go test -cover ./...

testv:
	go test -v -cover ./...

race:
	go test -race ./...

vet:
	go vet ./...

build:
	go build -o bin/pilad ./pilad

clean:
	rm -rf bin/ coverage.txt

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
