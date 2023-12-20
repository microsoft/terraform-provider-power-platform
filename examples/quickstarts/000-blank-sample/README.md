<!-- This document is auto-generated. Do not edit directly. Make changes to README.md.tmpl instead. -->
# Blank Sample (000 level)

This Terraform module aims to provide a fully managed infrastructure that integrates Microsoft's Power Platform and Azure services. Utilizing both powerplatform and azurerm Terraform providers, this module encapsulates best practices and serves as a reference architecture for scalable, reliable, and manageable cloud infrastructure.

## Prerequisites

- Azure subscription
- Power Platform environment
- Service Principal or User Account with appropriate permissions

## Example Files

The example files can be found in `examples/quickstarts/000-blank-sample`

## Terraform Version Constraints:
* `>= 1.5`

## Provider Requirements:
* **azurerm (`hashicorp/azurerm`):** (any version)
* **powerplatform (`microsoft/power-platform`):** (any version)

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `client_id` | The username of the Power Platform API in user@domain format | string | `null` | true |
| `secret` | The password of the Power Platform API user | string | `null` | false |
| `tenant_id` | The AAD tenant id of service principal or user | string | `null` | true |


## Output Values

| Name | Description |
|------|-------------|
| `id` | Unique identifier of the environment |
| `name` | Display name of the environment |
| `url` | URL of the environment |



## Resources
* `powerplatform_environment.development` from `powerplatform`

## Data Sources
* `data.powerplatform_connectors.all_connectors` from `powerplatform`


## Usage

Include this module in your Terraform scripts as follows:

```hcl

module "power_azure_infra" {
  source            = "./modules/000-blank-sample"
  resource_group    = "myResourceGroup"
  power_environment = "myPowerEnvironment"
}

```

## Detailed Behavior

### Azure Kubernetes Service

This module creates an AKS cluster to host containerized applications. The cluster is connected to the Azure Virtual Network.

### Azure SQL Database

The SQL database is provisioned with geo-redundant backups and is accessible only from within the Virtual Network.

### Power Platform Resources

A PowerApps Environment is created, tied to the Azure SQL Database for data storage. Power Automate flows are also defined and triggered by database changes, which in turn can run Kubernetes jobs.

### Network

All resources are provisioned within the same Azure Virtual Network, ensuring that they can communicate securely without exposure to the public internet.

## Limitations and Considerations

- Due to Power Platform limitations, certain resources may not fully support Terraform's state management.
- Make sure to set appropriate RBAC for Azure and Power Platform resources.
- This module is provided as a sample only and is not intended for production use without further customization.

## Additional Resources

- [Power Platform Admin Documentation](https://learn.microsoft.com/en-us/power-platform/admin/)
