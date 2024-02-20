### Using Azure CLI (Preferred)

The Power Platform provider can use the Azure CLI to authenticate. If you have the Azure CLI installed, you can use it to log in to your Azure account and the Power Platform provider will use the credentials from the Azure CLI.

## Creating Service Principal

You can follow this [guide](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app#register-an-application) to create a service principal.

## API Permissions

Following API permissions are required to use the Terraform Power Platform provider:

- Dynamics CRM
  - Dynamics CRM user_impersonation

- Power Platform API
  - AppManagement.ApplicationPackages.Install
  - AppManagement.ApplicationPackages.Read
  - Licensing.BillingPolicies.Read
  - Licensing.BillingPolicies.ReadWrite

- PowerApps Service
  - User

Or you can add them directly into your App Registration manifest:

```json
"requiredResourceAccess": [
  {
   "resourceAppId": "8578e004-a5c6-46e7-913e-12f58912df43",
   "resourceAccess": [
    {
     "id": "61bfce59-bddc-493f-b20c-32af5e904b83",
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
```

## Expose API

In "Expose API" menu of your App Registration, you need to define your application ID URI:

- Application ID URI: `api://<client_id>`, for example:

```bash
api://powerplatform_provider_terraform
```

- Add required scope:

1. Scope Name: `user_impersonation`
1. Who can consent: `Admins and users`
1. Admin consent display name: `Work with Power Platform Terraform Provider`
1. Admin consent description: `Allows connection to backend services of Power Platform Terraform Provider`
1. User consent display name: `Work with Power Platform Terraform Provider`
1. User consent description: `Allows connection to backend services of Power Platform Terraform Provider`
1. State: `Enabled`

Or you can add them directly into your App Registration manifest:

```json
 "oauth2Permissions": [
  {
   "adminConsentDescription": "Allows connection to backend services of Power Platform Terraform Provider",
   "adminConsentDisplayName": "Work with Power Platform Terraform Provider",
   "id": "2aedce72-ddc7-431d-920c-a321297ffdc2",
   "isEnabled": true,
   "lang": null,
   "origin": "Application",
   "type": "User",
   "userConsentDescription": "Allows connection to backend services of Power Platform Terraform Provider",
   "userConsentDisplayName": "Work with Power Platform Terraform Provider",
   "value": "user_impersonation"
  }
 ],
```

- You will finially need to preuthorize Azure CLI to access your API by adding client application `04b07795-8ddb-461a-bbee-02f9e1bf7b46`

Or you can add them directly into your App Registration manifest:

```json "preAuthorizedApplications": [
  {
   "appId": "04b07795-8ddb-461a-bbee-02f9e1bf7b46",
   "permissionIds": [
    "2aedce72-ddc7-431d-920c-a321297ffdc2"
   ]
  }
 ]
```

## Usage

After above steps you should be able to authenticate using Azure CLI:

```bash
az login  --allow-no-subscriptions --scope api://powerplatform_provider_terraform/.default
```
