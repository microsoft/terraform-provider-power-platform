#!/bin/sh
set -eux

echo "Installing tools need for GitHub Agent"

VERSION=1.21.1
curl -LO https://github.com/miniscruff/changie/releases/download/v${VERSION}/changie_${VERSION}_linux_amd64.tar.gz
tar -xzf changie_${VERSION}_linux_amd64.tar.gz changie
mv changie /usr/local/bin/
rm changie_${VERSION}_linux_amd64.tar.gz

VERSION=2.0.1
curl -LO https://github.com/golangci/golangci-lint/releases/download/v${VERSION}/golangci-lint-${VERSION}-linux-amd64.tar.gz
tar -xzf golangci-lint-${VERSION}-linux-amd64.tar.gz golangci-lint-${VERSION}-linux-amd64/golangci-lint
mv golangci-lint-${VERSION}-linux-amd64/golangci-lint /usr/local/bin/golangci-lint
rm -rf golangci-lint-${VERSION}-linux-amd64
rm golangci-lint-${VERSION}-linux-amd64.tar.gz

VERSION=1.11.4 
OS=linux
ARCH=amd64
curl -LO https://releases.hashicorp.com/terraform/${VERSION}/terraform_${VERSION}_${OS}_${ARCH}.zip
unzip terraform_${VERSION}_${OS}_${ARCH}.zip terraform
mv terraform /usr/local/bin/
rm terraform_${VERSION}_${OS}_${ARCH}.zip

VERSION=0.21.0
OS=linux
ARCH=amd64
curl -LO https://github.com/hashicorp/terraform-plugin-docs/releases/download/v${VERSION}/tfplugindocs_${VERSION}_${OS}_${ARCH}.zip
unzip tfplugindocs_${VERSION}_${OS}_${ARCH}.zip tfplugindocs
mv tfplugindocs /usr/local/bin/
rm tfplugindocs_${VERSION}_${OS}_${ARCH}.zip

