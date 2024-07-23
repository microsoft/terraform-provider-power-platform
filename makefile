deps:
	go mod tidy

build:
	$(MAKE) deps
	go build -o ./bin/

install:
	$(MAKE) build
	go install

clean:
	go clean -testcache
	rm -rf ./bin
	rm -rf /go/bin/terraform-provider-power-platform

userdocs:
	tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"

servedocs:
	$(MAKE) userdocs
	mkdocs serve

unittest:
	export TF_ACC=0
	$(MAKE) clean
	$(MAKE) install
	go test -v ./... -run "^TestUnit"

acctest:
	export TF_ACC=1
	$(MAKE) clean
	$(MAKE) install
	go test -timeout 120m -v ./... -run "^TestAcc"

test:
	export TF_ACC=1
	$(MAKE) clean
	$(MAKE) install
	go test -timeout 120m -v ./...

lint:
	golangci-lint run
