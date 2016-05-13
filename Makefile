.PHONY: vet lint

default: vet get

all: vet get test

get:
	go get ./...

test:
	go test -cover ./...

testv:
	go test -v -cover ./...

vet:
	go list ./... | grep -v /vendor/ | xargs -L1 go vet

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint

pilad:	get
	$(GOPATH)/bin/pilad

gox:	get
	gox -output "dist/{{.OS}}/{{.Arch}}/$(git rev-parse HEAD)/{{.Dir}}" ./pilad

release:
	docker run --rm --name="piladb_release" -v "$(PWD)":/gopath/src/github.com/fern4lvarez/piladb -w /gopath/src/github.com/fern4lvarez/piladb tcnksm/gox:1.5 make gox
