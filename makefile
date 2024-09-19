BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

deps:
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
	go clean -testcache
	rm -rf ./bin
	rm -rf /go/bin/terraform-provider-power-platform

userdocs:
	go generate

unittest:
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=0 go test -p 16 -timeout 10m -v ./... -run "^TestUnit$(TEST)"

acctest:
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"

test:
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=1 go test -p 10 -timeout 300m -v ./...

lint:
	golangci-lint run
