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
	tfplugindocs validate --provider-name powerplatform

userstory:
	./scripts/user_story_prompt.sh

unittest:
	clear
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=0 go test -p 16 -timeout 10m -v -cover ./... -run "^TestUnit$(TEST)"

acctest:
	clear
	$(MAKE) clean
	$(MAKE) install
ifeq ($(USE_PROXY),1)
	HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
else
	TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
endif

test:
	clear
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=1 go test -p 10 -timeout 300m -v -cover ./...

coverage:
	clear
	$(MAKE) clean
	$(MAKE) install
	@echo "Changed files:"
	@gh pr diff --name-only
	@echo "Running tests"
	TF_ACC=0 go test -p 16 -timeout 10m -v -cover -coverprofile=test-coverage.out ./... -run "^TestUnit$(TEST)"
	@echo "Generating coverage report"
	go tool cover -func=test-coverage.out

netdump_aaa:
	for pid in $$(cat /proc/[0-9]*/status 2>/dev/null | grep -E "^Name:|^Pid:" | grep -A1 "^Name:[[:space:]]*mitmdump$$" | grep "^Pid:" | awk '{print $$2}'); do kill -9 $$pid; echo "Killed process $$pid"; done && \
	rm -f mitmproxy.dump && \
	mitmdump -p 8080 -w mitmproxy.dump

netdump:
	echo "asdasd"

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
