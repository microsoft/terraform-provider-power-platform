---
page_title: "Authenticating to Power Platform using the Azure Developer CLI"
subcategory: "Authentication"
description: |-
  <no value>
---

# Authenticating to Power Platform using the Azure Developer CLI

The Power Platform provider can use the [Azure Developer CLI (azd)](https://learn.microsoft.com/azure/developer/azure-developer-cli/) to authenticate to Power Platform services. If you have the Azure Developer CLI installed, you can use it to log in to your Microsoft Entra Id account and the Power Platform provider will use the credentials from the Azure Developer CLI.

## Prerequisites

1. [Install the Azure Developer CLI](https://learn.microsoft.com/azure/developer/azure-developer-cli/install-azd)
1. [Create an app registration for the Power Platform Terraform Provider](app_registration.md)
1. Login using the Azure Developer CLI

    ```bash
    azd auth login
    ```

    Configure the provider to use the Azure Developer CLI with the following code:

    ```terraform
    provider "powerplatform" {
      use_dev_cli = true
    }
    ```
