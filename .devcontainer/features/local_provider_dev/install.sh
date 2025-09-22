#!/bin/sh
set -eux

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

# Install mitmproxy for inspecting HTTP traffic
MITM_VERSION=11.1.3
wget https://downloads.mitmproxy.org/${MITM_VERSION}/mitmproxy-${MITM_VERSION}-linux-x86_64.tar.gz
tar -xf mitmproxy-${MITM_VERSION}-linux-x86_64.tar.gz
mv mitmproxy mitmdump mitmweb /usr/local/bin/
rm mitmproxy-${MITM_VERSION}-linux-x86_64.tar.gz

mitmproxy --version
echo "Start mitmdump so that it generates a client certificate, and kill it after 5 seconds with a successful return code"
timeout 5s mitmdump -p 8080 || true

install -D ~/.mitmproxy/mitmproxy-ca-cert.pem /usr/local/share/ca-certificates/mitmproxy-ca.crt
cp /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt.bak
sh -c 'cat ~/.mitmproxy/mitmproxy-ca-cert.pem >> /etc/ssl/certs/ca-certificates.crt'

tfenv install latest

go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0
# Removing the golangci-lint binary from GOROOT/bin to avoid conflicts duplicated binaries
rm -f "$(go env GOROOT)/bin/golangci-lint" || true

# Turn off telemetry for az cli
# https://github.com/Azure/azure-cli?tab=readme-ov-file#telemetry-configuration
az config set core.collect_telemetry=false

curl -fsSL https://aka.ms/install-azd.sh | bash
# Turn off telemetry for azd
export AZURE_DEV_COLLECT_TELEMETRY=no
