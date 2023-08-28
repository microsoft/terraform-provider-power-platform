#!/bin/sh
set -e

echo "Setting up local provider install"

cp terraform.rc /go/bin/

# install tfplugindocs
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# fix permissions issue running as vscode and not root
sudo chmod -R 777 /go/pkg/mod
mkdir /workspaces/terraform-provider-power-platform/ && sudo chown -R vscode /workspaces/terraform-provider-power-platform/

# install mkdocs
sudo apt update && sudo apt install -y python3-pip
pip3 install mkdocs

# install mockgen
go install github.com/golang/mock/mockgen@latest
