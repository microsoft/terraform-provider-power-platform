---
page_title: "Authenticating to Power Platform using the Azure Developer CLI"
subcategory: "Authentication"
description: |-
  <no value>
---

# Authenticating to Power Platform using the Azure Developer CLI

The Power Platform provider can use the [Azure Developer CLI (azd)](https://learn.microsoft.com/azure/developer/azure-developer-cli/) to authenticate to Power Platform services. If you have the Azure Developer CLI installed, you can use it to log in to your Microsoft Entra Id account and the Power Platform provider will use the credentials from the Azure Developer CLI.

## User Authentication

For user-based authentication, you only need to install the Azure Developer CLI and log in with your user account.

### Prerequisites

1. [Install the Azure Developer CLI](https://learn.microsoft.com/azure/developer/azure-developer-cli/install-azd)

### Authentication Steps

1. Login using the Azure Developer CLI:

    ```bash
    azd auth login
    ```

2. Configure the provider to use the Azure Developer CLI:

    ```terraform
    provider "powerplatform" {
      use_dev_cli = true
    }
    ```

## Service Principal Authentication

For service principal authentication using Azure Developer CLI, you need to create an app registration first.

### Prerequisites

1. [Install the Azure Developer CLI](https://learn.microsoft.com/azure/developer/azure-developer-cli/install-azd)
2. [Create an app registration for the Power Platform Terraform Provider](app_registration.md)

### Authentication Steps

1. Login using a Service Principal:

    ```bash
    azd auth login --client-id <CLIENT_ID> --client-secret <CLIENT_SECRET> --tenant-id <TENANT_ID>
    ```

2. Configure the provider to use the Azure Developer CLI:

    ```terraform
    provider "powerplatform" {
      use_dev_cli = true
    }
    ```
