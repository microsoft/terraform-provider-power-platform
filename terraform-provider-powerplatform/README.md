# Terraform Power Platform Provider

This is a Terraform provider for the Power Platform. It is currently in development and is not ready for production use. It is not yet published to the Terraform registry. You can build it locally and use it in your Terraform configuration.

See the [example](./example) directory for an example of how to use it.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.18
- choco install make
- <https://marketplace.visualstudio.com/items?itemName=golang.go>

## Building The Provider

```sh
go mod init power-platform-terraform-provider
go mod tidy
make build
make install
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

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
make testacc
```

## End User Guidance

### Pre-requisites  

Install:

1. [Terraform](https://www.terraform.io/downloads.html)
2. [Docker](https://www.docker.com/products/docker-desktop/)
3. Setup Terraform user account for Power Platform

You can create dedicated user account or use your admin account. This account will be used by terraform to access Power Platform as Global Administrator:

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
  username = "var.username"
  password = "var.password"
  host = "http://localhost:8080"
}
```

4. Setup Terraform service principal for Azure

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

### Setup

1. Download Docker Image

**Note**
Request Azure Container Registry credentials from provider's owner

Download the docker image:

```powershell
docker login myregistry.azurecr.io -u user -p your_password
docker pull myregisry.azurecr.io/terraform_api:latest
```

Run docker image:

```powershell
docker run -dt -p 8080:80 --name terraform_api myregistry.azurecr.io/terraform_api:latest
```

2. Setup custom terraform provider

Option 1 - building from source:

```powershell
cd /terraform-provider-powerplatform
make build && make install
```

Option 2 - coping manually

Copy the 'terraform-provider-powerplatform.exe' to the '%appdata%\terraform.d\plugins\registry.terraform.io\microsoft\powerplatform\0.2\windows_amd64' folder
