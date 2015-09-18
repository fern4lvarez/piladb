default: vet get

all: vet get test

get:
	go get ./...

test:
	go test -cover ./...

testv:
	go test -v -cover ./...
vet:
	go vet ./...

server: get
	$(GOPATH)/bin/pilad
