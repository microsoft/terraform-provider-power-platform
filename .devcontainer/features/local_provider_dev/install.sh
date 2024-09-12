#!/bin/sh
set -e

echo "Setting up local provider install"

# install tfplugindocs
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# install delve debugger
go install -v github.com/go-delve/delve/cmd/dlv@latest
