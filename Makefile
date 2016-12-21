.PHONY: vet lint

default: vet get test

all: vet lint get test

get:
	go get ./...

test:
	go list ./... | grep -v /vendor/ | xargs -L1 go test -cover -coverprofile=coverage.txt -covermode=atomic

testv:
	go list ./... | grep -v /vendor/ | xargs -L1 go test -v -cover -coverprofile=coverage.txt -covermode=atomic

vet:
	go list ./... | grep -v /vendor/ | xargs -L1 go vet

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint

pilad:	get
	$(GOPATH)/bin/pilad

gox:	get
	gox -output "dist/{{.OS}}/{{.Arch}}/$(git rev-parse HEAD)/{{.Dir}}" ./pilad

release:
	docker run --rm --name="piladb_release" -v "$(PWD)":/gopath/src/github.com/fern4lvarez/piladb -w /gopath/src/github.com/fern4lvarez/piladb tcnksm/gox:latest make gox
