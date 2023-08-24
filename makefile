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
	TF_ACC=0 go test -v ./...

acctest: 
	TF_ACC=1 go test -v ./...

test:
	TF_ACC=0 go test -v ./...
	TF_ACC=1 go test -v ./...

deps:
	go mod tidy
