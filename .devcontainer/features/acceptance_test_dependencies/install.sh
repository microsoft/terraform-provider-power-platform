#!/bin/sh
set -e

echo "Setting up acceptance test dependencies"

# use terraform to download all the providers that we use in acceptance tests
terraform init

# copy the providers to the /go/bin directory so that they are available to the acceptance tests
find .terraform/providers/ -type f -exec cp {} /go/bin \;

# remove the .terraform directory
rm .terraform.lock.hcl
rm -rf .terraform
