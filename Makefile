clean:
	@go clean -testcache
	@rm -rf ./bin
	@rm -rf /go/bin/terraform-provider-power-platform

deps: clean
	@go mod tidy

build: deps
	@go fmt ./...
	@go build -o ./bin/

install: deps
	@go fmt ./...
	@go install

userdocs:
	@go generate
	@tfplugindocs validate --provider-name powerplatform

userstory:
	@./scripts/user_story_prompt.sh

unittest: clean install
	@TF_ACC=0 go test -p 16 -timeout 10m -v -cover ./... -run "^TestUnit$(TEST)"

acctest: clean install
ifeq ($(USE_PROXY),1)
	@HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
else
	@TF_ACC=1 go test -p 10 -timeout 300m -v ./... -run "^TestAcc$(TEST)"
endif

test: clean install
	@TF_ACC=1 go test -p 10 -timeout 300m -v -cover ./...

coverage: clean install
	@echo "Changed files:"
	@gh pr diff --name-only || true
	@echo "Running tests"
	@TF_ACC=0 go test -p 16 -timeout 10m -v -cover -coverprofile=test-coverage.out ./... -run "^TestUnit$(TEST)"
	@echo "Generating coverage report"
	@go tool cover -func=test-coverage.out

netdump:
	@mitmdump -p 8080 -w /tmp/mitmproxy.dump

lint:
	@golangci-lint run

precommit: clean build lint unittest userdocs
