#!/bin/sh
set -eux

# Match the installed Go to the version go.mod requires, the same way CI does with
# actions/setup-go's go-version-file. When they differ GOTOOLCHAIN downloads the required
# toolchain but the base image's older tools stay on the build path, and any coverage build
# of a package without test files fails with "compile: version ... does not match go tool
# version", which breaks `make unittest`, `make coverage` and `make precommit`.
GO_REQUIRED_VERSION=$(awk '/^go [0-9]+\.[0-9]+/ { print $2; exit }' go.mod)
GO_INSTALLED_VERSION=$(sed -n '1s/^go//p' /usr/local/go/VERSION 2>/dev/null || echo "")
if [ "$GO_REQUIRED_VERSION" != "$GO_INSTALLED_VERSION" ]; then
    echo "Replacing Go ${GO_INSTALLED_VERSION:-<none>} with Go $GO_REQUIRED_VERSION required by go.mod..."
    GO_ARCH=$(go env GOARCH)
    curl -fsSL "https://go.dev/dl/go${GO_REQUIRED_VERSION}.linux-${GO_ARCH}.tar.gz" -o /tmp/go.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
fi
go version

# Turn commits auto sign to auto
# https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits
#git config commit.gpgsign true

#Install codeql extension
echo "Installing GitHub CodeQL CLI extension..."
gh extensions install github/gh-codeql

# Set git to use LF line endings on commit
git config --global core.autocrlf input

# Turn off telemetry for az cli (only if az is installed)
# https://github.com/Azure/azure-cli?tab=readme-ov-file#telemetry-configuration
if command -v az > /dev/null 2>&1; then
    az config set core.collect_telemetry=false
fi
