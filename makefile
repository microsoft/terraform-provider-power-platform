BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

deps:
	go version
	go mod tidy

build:
	$(MAKE) deps
	go fmt ./...
	go build -o ./bin/

install:
	$(MAKE) deps
	go fmt ./...
	go install

clean:
	go version
	go clean -testcache
	rm -rf ./bin
	rm -rf /go/bin/terraform-provider-power-platform

userdocs:
	go generate

unittest:
	clear
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=0 go test -p 16 -timeout 10m -v ./... -run "^TestUnit$(TEST)"

acctest:
	clear
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"

test:
	clear
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=1 go test -p 10 -timeout 300m -v ./...

lint:
	clear
	golangci-lint --version
	golangci-lint run

precommit:
	$(MAKE) clean
	$(MAKE) build
	$(MAKE) lint
	$(MAKE) unittest
	$(MAKE) userdocs
