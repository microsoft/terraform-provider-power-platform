#!/bin/sh
set -eux

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
