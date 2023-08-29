build:
	go mod tidy
	go build -o ./bin/

install:
	go mod tidy
	go install

userdocs:
	tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"

servedocs:
	tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"
	mkdocs serve

unittest:
	export TF_ACC=0
	go test -v ./... -run "^TestUnit"

acctest: 
	export TF_ACC=1
	go test -v ./... -run "^TestAcc"

test:
	export TF_ACC=1
	go test -v ./...

deps:
	go mod tidy
