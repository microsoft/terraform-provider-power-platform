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
  - Licensing.Allocations.ReadWrite
  - Licensing.BillingPolicies.ReadWrite
  - PowerApps.Apps.Play
  - PowerApps.Apps.Read
  - EnvironmentManagement.Environments.Read
  - EnvironmentManagement.Groups.ReadWrite
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

### Authorizing client applications

You will need to preauthorize Azure CLI to access your API by adding client application `04b07795-8ddb-461a-bbee-02f9e1bf7b46`
