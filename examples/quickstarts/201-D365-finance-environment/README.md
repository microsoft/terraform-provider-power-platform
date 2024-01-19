<!-- This document is auto-generated. Do not edit directly. Make changes to README.md.tmpl instead. -->
# D365 Finance Deployment (201 level)

This Terraform module aims to provide a template for automating and standardizing the deployment and management of D365 Finance environments. It utilizes the deployment model outlined at https://learn.microsoft.com/en-us/power-platform/admin/unified-experience/finance-operations-apps-overview .

## Prerequisites

- Service Principal or User Account with permissions configured as referenced in [this provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication) .
- A properly assigned D365 license (for example, a Dynamics 365 Finance or Dynamics 365 Supply Chain Management license). For more information on the new license requirements, see https://learn.microsoft.com/en-us/power-platform/admin/unified-experience/finance-operations-apps-overview#transition-from-an-environment-slot-purchasing-model-to-a-capacity-based-model .
- At least 1 gigabyte of available Operations and Dataverse database capacities

## Example Files

The example files can be found in `examples/quickstarts/201-D365-finance-environment`

## Terraform Version Constraints:
* `>= 1.5`

## Provider Requirements:
* **powerplatform (`microsoft/power-platform`):** (any version)

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `client_id` | The client ID of the service principal (app registration) | string | `null` | true |
| `currency_code` | The desired Currency Code for the environment, such as 'USD' | string | `null` | true |
| `d365_finance_environment_name` | The name of the D365 Finance environment, such as 'd365fin-dev1 | string | `null` | true |
| `domain` | The domain of the environment, such as 'd365fin-dev1' | string | `null` | true |
| `environment_type` | The type of environment to deploy, such as 'Sandbox' | string | `null` | true |
| `language_code` | The desired Language Code for the environment, such as '1033' (U.S. english) | string | `null` | true |
| `location` | The region where the environment will be deployed, such as 'unitedstates' | string | `null` | true |
| `secret` | The client secret of the service principal (app registration) | string | `null` | true |
| `security_group_id` | The security group the environment will be associated with, a GUID. Can be set to 00000000-0000-0000-0000-000000000000 to indicate no security group restricting Dataverse access. | string | `null` | true |
| `tenant_id` | The Entra (AAD) tenant id of service principal or user | string | `null` | true |


## Output Values

| Name | Description |
|------|-------------|
| `id` | Unique identifier of the environment |
| `linked_app_id` | Unique identifier of the linked D365 Finance app |
| `linked_app_type` | Type of the linked D365 app |
| `linked_app_url` | URL of the linked D365 Finance app |
| `name` | Display name of the environment |
| `url` | URL of the environment |



## Resources
* `powerplatform_environment.xpp-dev1` from `powerplatform`


## Usage

Include this module in your Terraform scripts as follows:

```hcl

module "d365_finance_environment" {
  source            = "./modules/201-D365-finance-environment"
  client_id = "Your App Registration ID (GUID) here"
  secret = "Your App Registration Secret here"
  tenant_id = "Your Entra (Azure) Tenant ID (GUID) here"
  d365_finance_environment_name = "d365fin-dev1"
  location = "unitedstates"
  language_code = "1033"
  currency_code = "USD"
  environment_type = "Sandbox"
  security_group_id = "00000000-0000-0000-0000-000000000000"
  domain = "d365fin-dev1"
}

```

## Detailed Behavior

### Power Platform Environment

This module creates a Power Platform environment using a combination of the parameters in the terraform files as well as the default settings specified by the 'templates' property.

### Dynamics 365 Finance Environment

This module creates a Dynamics 365 Finance development environment using the default settings specified by the 'templates' and 'template_metadata' properties.

## Limitations and Considerations

- Provisioning can take over an hour, so refrain from rerunning the same environment creation Terraform files more than hourly, as this will cause unexpected behavior.
- This quickstart is configured for service-principal-based authentication as outlined in [this provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication). If you plan to use user-based authentication, you will need to ensure that the selected user is assigned a D365 Finance or D365 Supply Chain Management license as outlined in the [Unified Admin Experience Overview](https://learn.microsoft.com/en-us/power-platform/admin/unified-experience/finance-operations-apps-overview).

## Additional Resources

- [Power Platform Admin Documentation](https://learn.microsoft.com/en-us/power-platform/admin/)
