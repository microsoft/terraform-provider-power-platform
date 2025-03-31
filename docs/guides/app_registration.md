---
page_title: "Authentication: Creating an App Registration to use the Power Platform Provider"
subcategory: "Authentication"
description: |-
  <no value>
---

# Creating an App Registration to use the Power Platform Provider

The following steps will guide you through the process of manually creating an App Registration in Azure Active Directory to use the Power Platform Provider, but if you would like a script to run, see the [bootstrap scripts in the Power Platform QuickStarts](https://github.com/microsoft/power-platform-terraform-quickstarts/blob/main/bootstrap/tenant-configuration/main.tf)

## Register an Application

[The basics of how to create an app registration in Entra](https://learn.microsoft.com/entra/identity-platform/quickstart-register-app#register-an-application) are covered in Entra documentation.  Familiarize yourself with the process and then follow the steps below to create an app registration for the Power Platform Provider.

## API Permissions

Following API permissions are required to use the Terraform Power Platform provider:

- Dynamics CRM
  - Dynamics CRM user_impersonation

- Power Platform API
  - AppManagement.ApplicationPackages.Install
  - AppManagement.ApplicationPackages.Read
  - Licensing.Allocations.Read
  - Licensing.Allocations.ReadWrite
  - Licensing.BillingPolicies.Read
  - Licensing.BillingPolicies.ReadWrite
  - PowerApps.Apps.Play
  - PowerApps.Apps.Read
  - EnvironmentManagement.Environments.Read
  - EnvironmentManagement.Groups.Read
  - EnvironmentManagement.Groups.ReadWrite
  - EnvironmentManagement.Settings.Read
  - EnvironmentManagement.Settings.ReadWrite

- PowerApps Service
  - User

!> Note: If you don't see Power Platform API showing up in the list when searching by GUID, it's possible that you still have access to it but the visibility isn't refreshed. To force a refresh run the below PowerShell script:

### API Permissions configuration with manifest.

Below is an example JSON configuration for the required API permissions. This can be used in the app registration manifest:

```json
{
  "requiredResourceAccess": [
    {
      "resourceAppId": "8578e004-a5c6-46e7-913e-12f58912df43",
      "resourceAccess": [
        {
          "id": "61bfce59-bddc-493f-b20c-32af5e904b83",
          "type": "Scope"
        },
        {
          "id": "f1a0b2d4-3c5e-4b8c-9f7d-6a0e1f3a2b8e",
          "type": "Scope"
        },
        {
          "id": "9dafb9c1-c236-48b1-b142-20dcaab58675",
          "type": "Scope"
        },
        {
          "id": "048eb363-c1da-41d5-9edf-423b605ff23e",
          "type": "Scope"
        },
        {
          "id": "73cf5c38-5257-4f28-8bbb-f78acf3290a4",
          "type": "Scope"
        },
        {
          "id": "25223ba4-e810-4f08-9803-cde4b2057a13",
          "type": "Scope"
        },
        {
          "id": "a8f422ae-8922-45d4-a8f1-275a6bd43077",
          "type": "Scope"
        },
        {
          "id": "adef0bc0-3a5b-457a-834c-cabd82f0a6d2",
          "type": "Scope"
        },
        {
          "id": "3f4998a4-cbb8-4e1e-9ea0-fd7fc110bb74",
          "type": "Scope"
        }
      ]
    },
    {
      "resourceAppId": "00000003-0000-0000-c000-000000000000",
      "resourceAccess": [
        {
          "id": "e1fe6dd8-ba31-4d61-89e7-88639da4683d",
          "type": "Scope"
        }
      ]
    },
    {
      "resourceAppId": "475226c6-020e-4fb2-8a90-7a972cbfc1d4",
      "resourceAccess": [
        {
          "id": "0eb56b90-a7b5-43b5-9402-8137a8083e90",
          "type": "Scope"
        }
      ]
    },
    {
      "resourceAppId": "00000007-0000-0000-c000-000000000000",
      "resourceAccess": [
        {
          "id": "78ce3f0f-a1ce-49c2-8cde-64b5c0896db4",
          "type": "Scope"
        }
      ]
    }
  ]
}


```powershell
#Install the Microsoft Entra the module
Install-Module AzureAD

Connect-AzureAD
New-AzureADServicePrincipal -AppId 8578e004-a5c6-46e7-913e-12f58912df43 -DisplayName "Power Platform API"
```

!> Note: The `resourceAppId` values are the application IDs of the services in the Public cloud.  If you are [using a sovereign cloud](./nonpublic_clouds.md) the IDs will be different and you will need to use the appropriate application IDs for those services.

## Expose API

In "Expose an API" menu of your App Registration, you need to define your application ID URI:

- Application ID URI: `api://<client_id>`, for example:

```plaintext
api://powerplatform_provider_terraform
```

### Define scopes

1. Scope Name: `access`
1. Who can consent: `Admins and users`
1. Admin consent display name: `Work with Power Platform Terraform Provider`
1. Admin consent description: `Allows connection to backend services of Power Platform Terraform Provider`
1. User consent display name: `Work with Power Platform Terraform Provider`
1. User consent description: `Allows connection to backend services of Power Platform Terraform Provider`
1. State: `Enabled`

### Authorizing client applications

You will finially need to preuthorize Azure CLI to access your API by adding client application `04b07795-8ddb-461a-bbee-02f9e1bf7b46`

## Usage

After above steps you should be able to authenticate using Azure CLI:

```bash
az login --scope api://powerplatform_provider_terraform/.default
```

If your tenant doesn't have any Azure subscriptions, you can use the `--allow-no-subscriptions` flag to login. If you are working in a web-based devcontainer and you need more control over the interactive login process you can use the `--use-device-code` flag.
