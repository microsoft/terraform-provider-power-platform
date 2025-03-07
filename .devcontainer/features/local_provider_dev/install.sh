#!/bin/sh
set -e

echo "Setting up local provider install"

# install tfplugindocs
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest

# install delve debugger
go install -v github.com/go-delve/delve/cmd/dlv@latest

# install changie
go install github.com/miniscruff/changie@latest

# Get the latest GitHub CLI version
echo "Fetching latest GitHub CLI release..."
LATEST_URL=$(curl -s -I https://github.com/cli/cli/releases/latest | grep -i location | cut -d' ' -f2 | tr -d '\r')
GH_VERSION=$(echo $LATEST_URL | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+$' | tr -d 'v')

echo "Installing GitHub CLI version $GH_VERSION..."
wget https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_amd64.tar.gz -O /tmp/ghcli.tgz

tar -xzf /tmp/ghcli.tgz -C /tmp
mv /tmp/gh_${GH_VERSION}_linux_amd64/bin/gh /usr/local/bin/
chmod +x /usr/local/bin/gh
rm -rf /tmp/gh_${GH_VERSION}_linux_amd64
rm /tmp/ghcli.tgz
