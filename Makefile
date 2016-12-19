.PHONY: vet lint

default: vet get test

all: ## Run vet and lint, compiles, and run tests
	vet lint get test

get: ## Compiles go and store binary into $GOPATH/bin
	go get ./...

test: ## Run tests
	go list ./... | grep -v /vendor/ | xargs -L1 go test -cover

testv: ## Run tests with verbosity
	go list ./... | grep -v /vendor/ | xargs -L1 go test -v -cover

vet: ## Run `go vet` against all subpackages
	go list ./... | grep -v /vendor/ | xargs -L1 go vet

lint: ## Run Go linter against all subpackages
	go list ./... | grep -v /vendor/ | xargs -L1 golint

pilad: ## Compiles code and starts `pilad` server
	get
	$(GOPATH)/bin/pilad

gox:	## Cross compiles `pilad` in all OSs and Architectures and store binaries into dist/
	get
	gox -output "dist/{{.OS}}/{{.Arch}}/$(git rev-parse HEAD)/{{.Dir}}" ./pilad

release: ## Cross compiles `pilad` inside a container and store binaries into dist/
	docker run --rm --name="piladb_release" -v "$(PWD)":/gopath/src/github.com/fern4lvarez/piladb -w /gopath/src/github.com/fern4lvarez/piladb tcnksm/gox:latest make gox

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
