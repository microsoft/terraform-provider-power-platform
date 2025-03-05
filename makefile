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

debug:
	$(MAKE) install
	/usr/bin/env GOPATH=/home/runtimeuser/go TF_ACC=true /home/runtimeuser/go/bin/dlv dap --client-addr=:35119
	# dlv exec $(GOPATH)/bin/terraform-provider-power-platform --headless --continue --listen=:2345 --api-version=2 --log --accept-multiclient

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
