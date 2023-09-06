---
page_title: "Provider: Power Platform"
description: |-
  The Power Platform Provider allows managing environments and other resources within [Power Platform](https://powerplatform.microsoft.com/)
---

# Power Platform Provider

The Power Platform provider allows managing environments and other resources within [Power Platform](https://powerplatform.microsoft.com/).

**⚠️ WARNING:** This code is experimental and provided solely for evaluation purposes. It is **NOT** intended for production use and may contain bugs, incomplete features, or other issues. Use at your own risk, as it may undergo significant changes without notice, and no guarantees or support are provided. By using this code, you acknowledge and agree to these conditions. Consult the documentation or contact the maintainer if you have questions or concerns.

## Requirements

This provider requires **Terraform >= 0.12**.  For more information on provider installation and constraining provider versions, see the [Provider Requirements documentation](https://developer.hashicorp.com/terraform/language/providers/requirements).

## Installation

**ℹ INFO:** This provider is not yet published to the Terraform registry, so it requires a local installation to use it at this time.

To use the provider you can download the binaries from [Releases](https://github.com/microsoft/terraform-provider-power-platform/releases) to your local file system and configure Terraform to use your local mirror.  See the [Explicit Installation Method Configuration](https://developer.hashicorp.com/terraform/cli/config/config-file#explicit-installation-method-configuration) for more information about using local binaries.

```terraform
provider_installation  {
  filesystem_mirror {
    path = "/usr/share/terraform/providers"
    include = ["registry.terraform.io/microsoft/power-platform"]
  }
}
```

## Authentication

The provider allows authentication via service principal or user credentials. All sensitive information should be passed into Terraform using environment variables (don't put secrets in your tf files).

### Using a Service Principal (Preferred)

To access Power Platform APIs using a service principal, you need to register a new service principal application in your own Azure Active Directory (Azure AD) tenant and then register that same application with Power Platform.

You can find more information on how to do this in the following articles:

- [Programmability and Extensibility - Authentication - Power Platform | Microsoft Learn](https://learn.microsoft.com/en-us/power-platform/admin/programmability-authentication-v2)
- [PowerShell: Create a service principal - Power Platform | Microsoft Learn](https://learn.microsoft.com/en-us/power-platform/admin/powershell-create-service-principal).
- [Registering an Admin Management Application](https://learn.microsoft.com/en-us/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application)

```terraform
# Configure the Power Platform Provider using a service principal
provider "powerplatform" {
  client_id = var.client_id
  secret    = var.secret
  tenant_id = var.tenant_id
}
```

```bash
export TF_VAR_client_id=<client_id>
export TF_VAR_secret=<secret>
export TF_VAR_tenant_id=<tenant_id>
```

### Using Username and Password

```terraform
# Configure the Power Platform Provider using username/password
provider "powerplatform" {
  username  = var.username
  password  = var.password
  tenant_id = var.tenant_id
}
```

```bash
export TF_VAR_username=<username>
export TF_VAR_password=<password>
export TF_VAR_tenant_id=<tenant_id>
```
###  Creating a "secret.tfvars" file to store your credentials

Alternatively you can create a "secret.tfvars" file and execute the "terraform plan" command specifying a local variables file:

```bash
# terraform plan command pointing to a secret.tfvars
terraform plan -var-file="secret.tfvars"
```
Below you will find an example of how to create your "secret.tfvars" file, remember to specify the correct path of it when executing.
We include "*.tfvars" in .gitignore to avoid save the secrets in it repository.

```bash
# sample "secret.tfvars" values
client_id = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
secret    = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
tenant_id = "XXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
```

In the terraform documentation ["Protect sensitive input variables"](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables) you can find more examples.

## Environment Variables

In addition to the variables that are passed into the provider, there are a few environment variables that can be used to configure the provider.

| Name | Description | Default Value |
|------|-------------|---------------|
| `POWER_PLATFORM_USERNAME` | The username to use for authentication. | |
| `POWER_PLATFORM_PASSWORD` | The password to use for authentication. | |
| `POWER_PLATFORM_CLIENT_ID` | The service principal client id | |
| `POWER_PLATFORM_SECRET` | The service principal secret | |
| `POWER_PLATFORM_TENANT_ID` | The guid of the tenant | |
| `POWER_PLATFORM_HOST` | The API endpoint used for managing Power Platform resources | |

Variables passed into the provider will override the environment variables.

## Resources and Data Sources

Use the navigation to the left to read about the available resources and data sources.

## Contributing

Contributions to this provider are always welcome! Please see the [Contribution Guidelines](https://github.com/microsoft/terraform-provider-power-platform/)
