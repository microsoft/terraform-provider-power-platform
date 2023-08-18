# Terraform Provider for Power Platform

The Power Platform Terraform Provider allows managing resources within Power Platform.

**⚠️ WARNING:** This code is experimental and provided solely for evaluation purposes. It is **NOT** intended for production use and may contain bugs, incomplete features, or other issues. Use at your own risk, as it may undergo significant changes without notice, and no guarantees or support are provided. By using this code, you acknowledge and agree to these conditions. Consult the documentation or contact the maintainer if you have questions or concerns.

**ℹ INFO:** This Power Platform provider is not yet published to the Terraform registry, but you can directly reference it from GitHub or build it locally and use it in your Terraform configuration.

## Requirements

## Usage Example

The following is a basic example of using the Terraform Provider for Power Platform.  See the [examples](./examples) directory for additional examples.

### Register the provider

```terraform
terraform {
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
    }
  }
}
```

### Configure the Power Platform Provider

The provider allows authentication via service principal or user credentials.  All sensitive information should be passed into Terraform using environment variables (don't put secrets in your tf files).

```terraform 
# Configure the Power Platform Provider using username/password
provider "powerplatform" {
  # The Power Platform Provider supports authenticating using Service Principal or user credentials.
  username = var.username
  password = var.password
  tenant_id = var.tenant_id
}
```

Or configure the Power Platform Provider using a service principal.  Service principals should be configured following [this](https://learn.microsoft.com/en-us/power-platform/admin/programmability-authentication-v2) and [this](https://learn.microsoft.com/en-us/power-platform/admin/powershell-create-service-principal).

```terraform 
# Configure the Power Platform Provider using username/password
provider "powerplatform" {
  # The Power Platform Provider supports authenticating using Service Principal or user credentials.
  client_id = var.client_id
  secret = var.secret
  tenant_id = var.tenant_id
}
```

### Create a Power Platform Environment

```terraform 
resource "powerplatform_environment" "development" {
  display_name     = "example_environment"
  location         = "europe"
  language_code    = "1033"
  currency_code    = "USD"
  environment_type = "Sandbox"
  domain           = "mydomain"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}
```

## Developing and Contributing to the Provider

The [DEVELOPER.md](DEVELOPER.md) file is a basic outline on how to build and develop the provider.

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit <https://cla.opensource.microsoft.com>.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft trademarks or logos is subject to and must follow [Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/legal/intellectualproperty/trademarks/usage/general). Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship. Any use of third-party trademarks or logos are subject to those third-party's policies.
