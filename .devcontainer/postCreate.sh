#!/bin/sh
set -eux

# Turn commits auto sign to auto
# https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits
#git config commit.gpgsign true

#Install codeql extension
echo "Installing GitHub CodeQL CLI extension..."
gh extensions install github/gh-codeql
git config --global core.autocrlf input
