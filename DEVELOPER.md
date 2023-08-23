# Developer Guide

## Developer Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.19
- [Visual Studio Code](https://code.visualstudio.com/)
- [Go Language Support for VS Code](https://marketplace.visualstudio.com/items?itemName=golang.go)

## Building The Provider

```sh
go mod tidy
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```sh
go get github.com/author/dependency
go mod tidy
```

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#developer-requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

To build run `go build`

To install locally run `go install .`

## End User Guidance

### Pre-requisites  

```json
terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  tenant_id = var.tenant_id
}
```

1. Setup Terraform service principal for Azure

Terraform will need a service principal to access Azure. List of required permissions can be found [here](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret)

```json
terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.15.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.0"
    }
  }
}

provider "azuread" {
  tenant_id     = "var.aad_tenant_id"
  client_id     = "var.aad_client_id"
  client_secret = "var.aad_client_secret"
}

provider "azurerm" {

  subscription_id = "var.azure_subscription_id"
  client_id       = "var.aad_client_id"
  client_secret   = "var.aad_client_secret"
  tenant_id       = "var.aad_tenant_id"
}
```

## Running Provider locally in VSCode (linux)

1. Open VSCode, Open with devcontainer. This will install go and terraform tools.
2. Open bash terminal
3. `cd` to the same parent folder as this ReadMe.
4. Execute commands below:

```bash

go install

export TF_CLI_CONFIG_FILE=<path to dev.tfrc> # dev.tfrc
export POWER_PLATFORM_USERNAME=<username>
export POWER_PLATFORM_PASSWORD=<password>
export POWER_PLATFORM_HOST=<api url>

# Navigate to a folder that contains main.tf and run below
terraform plan
```

Note: You cannot run `terraform init` when using dev overrides. `terraform init` will validate the versions and provider source, while `terraform plan` will skip those validations when `dev overrides` is part of your config.

## Debugging provider in VSCode

1. Open VSCode with the root folder as the parent of this ReadMe
1. Click On Run and Debug (F5)
1. Copy `TF_REATTACH_PROVIDERS` value in the Debug Console
1. Run `export TF_REATTACH_PROVIDERS=<value>` with the value copied from the above step
1. Add breakpoints
1. `cd` to a parent folder where main.tf exists
1. Run `terraform apply`

## Running Acceptance Tests

1. Set variable `TF_ACC` to `1` by running `export TF_ACC=1`
1. Go to provider's root folder and run `go test -v ./...`
1. To run single acceptance test run `go test -v ./... -run TestAcc<test_name>`

## Running Unit Tests

1. Validate that variable TF_ACC is not set by running `echo $TF_ACC`
1. Go to provider's root folder and run `go test -v ./...`
1. To run single unit test run `go test -v ./... -run TestUnit<test_name>`

## Generating Documentation

1. Execute below command to generate documentation

```bash
tfplugindocs generate --provider-name powerplatform --rendered-provider-name "Power Platform"
```

## Generating Mocks for Unit Tests

1. Execute below command to generate mocks

```bash
mockgen -destination=./internal/mocks/client_mocks_bapi.go -package=powerplatform_mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi" ApiClientInterface
```
