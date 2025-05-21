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
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=0 go test -p 16 -timeout 10m -v -cover ./... -run "^TestUnit$(TEST)"

acctest:
	$(MAKE) clean
	$(MAKE) install
ifeq ($(USE_PROXY),1)
	HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
else
	TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
endif

test:
	$(MAKE) clean
	$(MAKE) install
	TF_ACC=1 go test -p 10 -timeout 300m -v -cover ./...

coverage:
	$(MAKE) clean
	$(MAKE) install
	@echo "Changed files:"
	@gh pr diff --name-only
	@echo "Running tests"
	TF_ACC=0 go test -p 16 -timeout 10m -v -cover -coverprofile=test-coverage.out ./... -run "^TestUnit$(TEST)"
	@echo "Generating coverage report"
	go tool cover -func=test-coverage.out

netdump:
	mitmdump -p 8080 -w /tmp/mitmproxy.dump

lint:
	golangci-lint --version
	golangci-lint run

precommit:
	$(MAKE) clean
	$(MAKE) build
	$(MAKE) lint
	$(MAKE) unittest
	$(MAKE) userdocs

# This command is for copilot agent. Wehn using out devcontainer, you will have all the tools installed already.
installtools:
	OS=linux
	ARCH=amd64
	CHANGIE_VERSION=1.21.1
	LINTER_VERSION=2.0.1
	TERRAFORM_VERSION=1.11.4 
	TF_PLUGIN_DOCS_VERSION=0.21.0
	curl -LO https://github.com/miniscruff/changie/releases/download/v${CHANGIE_VERSION}/changie_${CHANGIE_VERSION}_linux_amd64.tar.gz
	tar -xzf changie_${CHANGIE_VERSION}_linux_amd64.tar.gz changie
	mv changie /usr/local/bin/
	rm changie_${CHANGIE_VERSION}_linux_amd64.tar.gz
	curl -LO https://github.com/golangci/golangci-lint/releases/download/v${LINTER_VERSION}/golangci-lint-${LINTER_VERSION}-linux-amd64.tar.gz
	tar -xzf golangci-lint-${LINTER_VERSION}-linux-amd64.tar.gz golangci-lint-${LINTER_VERSION}-linux-amd64/golangci-lint
	mv golangci-lint-${LINTER_VERSION}-linux-amd64/golangci-lint /usr/local/bin/golangci-lint
	rm -rf golangci-lint-${LINTER_VERSION}-linux-amd64
	rm golangci-lint-${LINTER_VERSION}-linux-amd64.tar.gz
	curl -LO https://releases.hashicorp.com/terraform/${VERSION}/terraform_${TERRAFORM_VERSION}_${OS}_${ARCH}.zip
	unzip terraform_${TERRAFORM_VERSION}_${OS}_${ARCH}.zip terraform
	mv terraform /usr/local/bin/
	rm terraform_${TERRAFORM_VERSION}_${OS}_${ARCH}.zip
	curl -LO https://github.com/hashicorp/terraform-plugin-docs/releases/download/v${TF_PLUGIN_DOCS_VERSION}/tfplugindocs_${TF_PLUGIN_DOCS_VERSION}_${OS}_${ARCH}.zip
	unzip tfplugindocs_${TF_PLUGIN_DOCS_VERSION}_${OS}_${ARCH}.zip tfplugindocs
	mv tfplugindocs /usr/local/bin/
	rm tfplugindocs_${TF_PLUGIN_DOCS_VERSION}_${OS}_${ARCH}.zip
