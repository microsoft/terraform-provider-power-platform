---
page_title: "Authentication: Connecting to non-public sovereign clouds"
subcategory: "Authentication"
description: |-
  <no value>
---

# Connecting to non-public sovereign clouds

Microsoft offers sovereign clouds (for example [US Government](https://learn.microsoft.com/en-us/power-platform/admin/microsoft-dynamics-365-government) and [China](https://learn.microsoft.com/en-us/power-platform/admin/about-microsoft-cloud-china)) that are physically isolated instances of Power Platform, Entra ID, and Azure. These isolated clouds are designed to make sure that data residency, sovereignty, and compliance requirements are honored within geographical boundaries.

>! Warning: Sovereign clouds may not support the same features as the public cloud. Make sure to [check the documentation](https://aka.ms/bapfunctionalparity) for the specific cloud you are using.

## Use Case

This guide can be useful in case you need to connect the Power Platform Terraform Provider to one of these clouds.

## Procedure

Configure the Azure CLI to work with that Cloud:

```bash
az cloud set --name AzureUSGovernment | AzureChinaCloud
```

See the [Azure CLI documentation](https://learn.microsoft.com/en-us/cli/azure/manage-clouds-azure-cli) for more information on connecting to alternate clouds.

Login to the Azure CLI using:

```bash
az login --scope api://powerplatform_provider_terraform/.default
```

Configure Power Platform provider similar to the following example:

```terraform
provider "azurerm" {
  cloud = "china" # public (default) | gcc | gcchigh | dod | china | ex | rx
  ...
}
```

The `cloud` configuration parameter accepts `public` (default), `gcc`, `gcchigh`, `dod`, `china`, `ex`, or `rx`
