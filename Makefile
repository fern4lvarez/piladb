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

pilad: get
	$(GOPATH)/bin/pilad
