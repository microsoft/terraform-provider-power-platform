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
az login
```

If your tenant doesn't have any Azure subscriptions, you can use the `--allow-no-subscriptions` flag to login. If you are working in a web-based devcontainer and you need more control over the interactive login process you can use the `--use-device-code` flag.
