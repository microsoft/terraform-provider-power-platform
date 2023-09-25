<!-- This document is auto-generated. Do not edit directly. Make changes to README.md.tmpl instead. -->
# D365 Finance Deployment (201 level)

This Terraform module aims to provide a template for automating and standardizing the deployment and management of D365 Finance environments. It utilizes the deployment model outlined at https://learn.microsoft.com/en-us/power-platform/admin/unified-experience/finance-operations-apps-overview .

## Prerequisites

- Service Principal or User Account with appropriate permissions
- A properly assigned D365 license (for example, a Dynamics 365 Finance or Dynamics 365 Supply Chain Management license)
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
| `currency_code` | The desired Currency Code for the environment | string | `"USD"` | false |
| `d365_finance_environment_name` | The name of the D365 Finance environment | string | `"d365fin-environment"` | false |
| `domain` | The domain of the environment | string | `"sample-d365f-environment"` | false |
| `environment_type` | The type of environment to deploy | string | `"Sandbox"` | false |
| `language_code` | The desired Language Code for the environment | string | `"1033"` | false |
| `location` | The Azure region where the environment will be deployed | string | `"canada"` | false |
| `secret` | The client secret of the service principal (app registration) | string | `null` | true |
| `security_group_id` | The security group the environment will be associated with | string | `"00000000-0000-0000-0000-000000000000"` | false |
| `templates` | The list of application templates to use when deploying the environment. | list(string) | `["D365_FinOps_Finance"]` | false |
| `tenant_id` | The Entra (AAD) tenant id of service principal or user | string | `null` | true |


## Output Values

| Name | Description |
|------|-------------|
| `id` | Unique identifier of the environment |
| `name` | Display name of the environment |
| `url` | URL of the environment |



## Resources
* `powerplatform_environment.development` from `powerplatform`


## Usage

Include this module in your Terraform scripts as follows:

```hcl

module "d365_finance_environment" {
  source            = "./modules/201-D365-finance-environment"
}

```

## Detailed Behavior

### Power Platform Environment

This module creates a Power Platform environment using a combination of the parameters in the terraform files as well as the default settings specified by the 'templates' property.

### Dynamics 365 Finance Environment

This module creates a Dynamics 365 Finance environment using the default settings specified by the 'templates' and 'template_metadata' properties.

## Limitations and Considerations

- Provisioning can take over an hour, so refrain from rerunning the same environment creation Terraform files more than hourly, as this will cause unexpected behavior.
- Be sure the relevant users are assigned the correct Dynamics 365 licenses, as this can cause unexpected errors.

## Additional Resources

- [Power Platform Admin Documentation](https://learn.microsoft.com/en-us/power-platform/admin/)
