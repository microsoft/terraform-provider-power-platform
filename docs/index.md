---
page_title: "Provider: Power Platform"
description: |-
  The Power Platform Provider allows managing environments and other resources within [Power Platform](https://powerplatform.microsoft.com/)
---

# Power Platform Provider

The Power Platform provider allows managing environments and other resources within [Power Platform](https://powerplatform.microsoft.com/).

!> This code is made available as a public preview. Features are being actively developed and may have restricted or limited functionality. Future updates may introduce breaking changes, but we follow [Semantic Versioning](https://semver.org/) to help mitigate this. The software may contain bugs, errors, or other issues that could cause service interruption or data loss. We recommend backing up your data and testing in non-production environments. Your feedback is valuable to us, so please share any issues or suggestions you encounter via GitHub issues.

## Requirements

This provider requires **Terraform >= 0.12**.  For more information on provider installation and constraining provider versions, see the [Provider Requirements documentation](https://developer.hashicorp.com/terraform/language/providers/requirements).

## Installation

To use this provider, add the following to your Terraform configuration:

```terraform
terraform {
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
      version = "~> 3.1" # Replace with the latest version
    }
  }
}
```

See the official Terraform documentation for more information about [requiring providers](https://developer.hashicorp.com/terraform/language/providers/requirements).

## Authenticating to Power Platform

Terraform supports a number of different methods for authenticating to Power Platform.

* [Authenticating to Power Platform using the Azure CLI](#authenticating-to-power-platform-using-the-azure-cli)
* [Authenticating to Power Platform using a Service Principal with OIDC](#authenticating-to-power-platform-using-a-service-principal-with-oidc)
* [Authenticating to Power Platform using a Service Principal and a Client Secret](#authenticating-to-power-platform-using-a-service-principal-and-a-client-secret)
* [Authenticating to Power Platform using a Managed Identity](#authenticating-to-power-platform-using-a-managed-identity)
* [Authenticating to Power Platform using Azure DevOps Workload Identity Federation](#authenticating-to-power-platform-using-workload-identity-federation)

We recommend using either a Service Principal when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

Important Notes about Authenticating using the Azure CLI:

* Terraform only supports authenticating using the az CLI (and this must be available on your PATH) - authenticating using the older azure CLI or PowerShell Cmdlets are not supported.
* Authenticating via the Azure CLI is only supported when using a User Account. If you're using a Service Principal (for example via az login --service-principal) you should instead authenticate via the Service Principal directly (either using a Client Secret or OIDC).

### Authenticating to Power Platform using the Azure CLI

The Power Platform provider can use the [Azure CLI](https://learn.microsoft.com/cli/azure/) to authenticate to Power Platform services. If you have the Azure CLI installed, you can use it to log in to your Microsoft Entra Id account and the Power Platform provider will use the credentials from the Azure CLI.

#### Prerequisites

1. [Install the Azure CLI](https://docs.microsoft.com/cli/azure/install-azure-cli)
1. [Create an app registration for the Power Platform Terraform Provider](guides/app_registration.md)
1. Login using the scope as the "expose API" you configured when creating the app registration

    ```bash
    az login --allow-no-subscriptions --scope api://powerplatform_provider_terraform/.default
    ```

    Configure the provider to use the Azure CLI with the following code:

    ```terraform
    provider "powerplatform" {
      use_cli = true
    }
    ```

### Authenticating to Power Platform using a Service Principal with OIDC

The Power Platform provider can use a Service Principal with OpenID Connect (OIDC) to authenticate to Power Platform services. By using [Microsoft Entra's workload identity federation](https://learn.microsoft.com/entra/workload-id/workload-identity-federation), your CI/CD pipelines in GitHub or Azure DevOps can access Power Platform resources without needing to manage secrets.

1. [Create an app registration for the Power Platform Terraform Provider](guides/app_registration.md)
1. [Register your app registration with Power Platform](https://learn.microsoft.com/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application)
1. [Create a trust relationship between your CI/CD pipeline and the app registration](https://learn.microsoft.com/entra/workload-id/workload-identity-federation-create-trust?pivots=identity-wif-apps-methods-azp)
1. Configure the provider to use OIDC with the following code:

    ```terraform
    provider "powerplatform" {
      use_oidc = true
    }
    ```

Additional Resources about OIDC:

* [OpenID Connect authentication with Microsoft Entra ID](https://learn.microsoft.com/entra/architecture/auth-oidc)
* [Configuring OpenID Connect for GitHub and Microsoft Entra ID](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-azure)

### Authenticating to Power Platform using a Service Principal and a Client Secret

The Power Platform provider can use a Service Principal with Client Secret to authenticate to Power Platform services.

1. [Create an app registration for the Power Platform Terraform Provider](guides/app_registration.md)
1. [Register your app registration with Power Platform](https://learn.microsoft.com/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application)
1. Configure the provider to use a Service Principal with a Client Secret with either environment variables or using Terraform variables

### Authenticating to Power Platfomr using Service Principal and certificate

1. [Create an app registration for the Power Platform Terraform Provider](guides/app_registration.md)
1. [Register your app registration with Power Platform](https://learn.microsoft.com/power-platform/admin/powerplatform-api-create-service-principal#registering-an-admin-management-application)
1. Generate a certificate using openssl or other tools

    ```bash
    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365
    ```

1. Merge public and private part of the certificate files together

    Using linux shell

    ```bash
    cat *.pem > cert+key.pem
    ```

    Using Powershell

    ```powershell
    Get-Content .\cert.pem, .\key.pem | Set-Content cert+key.pem
    ```

1. Generate pkcs12 file

    ```bash
    openssl pkcs12 -export -out cert.pkcs12 -in cert+key.pem
    ```

1. Add public part of the certificate (`cert.pem` file) to the app registration
1. Store your key.pem and the password used to generate in a safe place
1. Configure the provider to use certificate with the following code:

    ```terraform
    provider "powerplatform" {
      client_id     = var.client_id
      tenant_id     = var.tenant_id
      client_certificate_file_path = "${path.cwd}/cert.pkcs12"
      client_certificate_password  = var.cert_pass
    }
    ```

### Authenticating to Power Platform Using a Managed Identity

The Power Platform provider can use a [Managed Identity](https://learn.microsoft.com/en-us/entra/identity/managed-identities-azure-resources/overview) (previously called Managed Service Identity, or MSI) to authenticate to Power Platform services for keyless authentication in scenarios where the provider is being executed in select Azure services, such as Microsoft-hosted or self-hosted Azure DevOps pipelines.

#### System-Managed Identity

1. [Enable system-managed identity on an Azure resource](https://learn.microsoft.com/en-us/entra/identity/managed-identities-azure-resources/overview)
1. Register the managed identity with the Power Platform using the Application ID from the enterprise application for the system-managed identity resource. This task can be performed using either [the Power Platform Terraform Provider itself](https://registry.terraform.io/providers/microsoft/power-platform/latest/docs/resources/admin_management_application), or [PowerShell]([Register the managed identity with the Power Platform](https://learn.microsoft.com/en-us/power-platform/admin/powershell-create-service-principal).
1. Configure the provider to use the system-managed identity. Note that no Client ID is required as the Client ID is derived from the Azure resource running the provider.

    ```terraform
    provider "powerplatform" {
      use_msi = true
    }
    ```

#### User-Managed Identity

1. [Create a User-Managed Identity resource](https://learn.microsoft.com/en-us/entra/identity/managed-identities-azure-resources/overview)
1. Register the managed identity with the Power Platform using the Application ID from the enterprise application for the system-managed identity resource. This task can be performed using either [the Power Platform Terraform Provider itself](https://registry.terraform.io/providers/microsoft/power-platform/latest/docs/resources/admin_management_application), or [PowerShell]([Register the managed identity with the Power Platform](https://learn.microsoft.com/en-us/power-platform/admin/powershell-create-service-principal).
1. Configure the provider to use the System-Managed Identity. Note that this example sets the Client ID in the provider configuration, but it could also be set using the POWER_PLATFORM_CLIENT_ID environment variable.

    ```terraform
    provider "powerplatform" {
      use_msi = true
      client_id = var.client_id # This should be the Client ID from the user-managed identity resource.
    }
    ```

### Authenticating to Power Platform Using Azure DevOps Workload Identity Federation

The Power Platform provider can use [Azure DevOps Workload Identity Federation](https://devblogs.microsoft.com/devops/introduction-to-azure-devops-workload-identity-federation-oidc-with-terraform/) with Azure DeOps pipelines to authenticate to Power Platform services.

*Note: For similar hands-off authentication in GitHub and Azure DevOps, the Power Platform Provider also supports the [OIDC authentication method](#authenticating-to-power-platform-using-a-service-principal-with-oidc).*

1. Create an [App Registration](guides/app_registration.md) or a [User-Managed Identity](https://learn.microsoft.com/en-us/entra/identity/managed-identities-azure-resources/overview). This resource will be used to manage the identity federation with Azure DevOps.
1. Register the App Registration or Managed Identity with the Power Platform. This task can be performed using [the provider itself](/resources/admin_management_application.md) or [PowerShell](https://learn.microsoft.com/en-us/power-platform/admin/powershell-create-service-principal).
1. [Complete the service connection configuration in Azure and Azure DevOps](https://learn.microsoft.com/en-us/azure/devops/pipelines/release/configure-workload-identity?view=azure-devops&tabs=managed-identity). Note that Azure DevOps may automatically generate the federated credential in Azure, depending on your permissions and Azure Subscription configuration.
1. Configure the provider to use Azure DevOps Workload Identity Federation. This authentication option also requires values to be set in the ARM_OIDC_REQUEST_TOKEN and POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID environment variables, which should be configured in the AzDO pipeline itself. Note that this example sets some of the required properties in the provider configuration, but the whole configuration could also be performed using just environment variables.

    ```terraform
    provider "powerplatform" {
      tenant_id = var.tenant_id
      client_id = var.client_id # The client ID for the Azure resource containing the federated credentials for Azure DevOps. Should be an App Registration or a Managed Identity.
    }
    ```

### Using Environment Variables

We recommend using Environment Variables to pass the credentials to the provider.

| Name | Description | Default Value |
|------|-------------|---------------|
| `POWER_PLATFORM_CLIENT_ID` | The service principal client id | |
| `POWER_PLATFORM_CLIENT_SECRET` | The service principal secret | |
| `POWER_PLATFORM_TENANT_ID` | The guid of the tenant | |
| `POWER_PLATFORM_CLOUD` | override for the cloud used (default is `public`) | |
| `POWER_PLATFORM_USE_OIDC` | if set to `true` then OIDC authentication will be used | |
| `POWER_PLATFORM_USE_CLI` | if set to `true` then Azure CLI authentication will be used | |
| `POWER_PLATFORM_USE_MSI` | if set to `true` then Managed Identity authentication will be used | |
| `POWER_PLATFORM_CLIENT_CERTIFICATE` | The Base64 format of your certificate that will be used for certificate-based authentication | |
| `POWER_PLATFORM_CLIENT_CERTIFICATE_FILE_PATH` | The path to the certificate that will be used for certificate-based authentication | |
| `POWER_PLATFORM_AZDO_SERVICE_CONNECTION_ID` | The GUID of the Azure DevOps service connection to be used for Azure DevOps Workload Identity Federation | |

-> Variables passed into the provider will override the environment variables.

#### Using Terraform Variables

Alternatively, you can configure the provider using variables in your Terraform configuration which can be passed in via [command line parameters](https://developer.hashicorp.com/terraform/language/values/variables#variables-on-the-command-line), [a `*.tfvars` file](https://developer.hashicorp.com/terraform/language/values/variables#variable-definitions-tfvars-files), or [environment variables](https://developer.hashicorp.com/terraform/language/values/variables#environment-variables).  If you choose to use variables, please be sure to [protect sensitive input variables](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables) so that you do not expose your credentials in your Terraform configuration.

```terraform
provider "powerplatform" {
  # Use a service principal to authenticate with the Power Platform service
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
}
```

## Additional configuration

In addition to the authentication options, the following options are also supported in the provider block:

| Name | Description | Default Value |
|------|-------------|---------------|
| `telemetry_optout` | Opting out of telemetry will remove the User-Agent and session id headers from the requests made to the Power Platform service.  There is no other telemetry data collected by the provider.  This may affect the ability to identify and troubleshoot issues with the provider. | `false` |


If you are using Azure CLI for authentication, you can also turn off CLI's telemetry by executing the following [command](https://github.com/Azure/azure-cli?tab=readme-ov-file#telemetry-configuration):
```bash 
az config set core.collect_telemetry=false
```

## Resources and Data Sources

Use the navigation to the left to read about the available resources and data sources.

!> By calling `terraform destroy` all the resources, that you've created, will be deleted permamently deleted. Please be careful with this command when working with production environments. You can use [prevent-destroy](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy) lifecycle argument in your resources to prevent accidental deletion.  

## Examples

More detailed examples can be found in the [Power Platform Terraform Quickstarts](https://github.com/microsoft/power-platform-terraform-quickstarts) repo.  This repo contains a number of examples for using the Power Platform provider to manage environments and other resources within Power Platform along with Azure and Entra.

## Releases

A full list of released versions of the Power Platform Terraform Provider can be found [here](https://github.com/microsoft/terraform-provider-power-platform/releases).  Starting from v3.0.0, a summary of the changes to the provider in each release are documented the [CHANGELOG.md file in the GitHub repository](https://github.com/microsoft/terraform-provider-power-platform/blob/main/CHANGELOG.md). This provider follows Semantic Versioning for releases. The provider version is incremented based on the type of changes included in the release.

## Contributing

Contributions to this provider are always welcome! Please see the [Contribution Guidelines](https://github.com/microsoft/terraform-provider-power-platform/)
