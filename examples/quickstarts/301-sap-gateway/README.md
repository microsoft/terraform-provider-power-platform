# Gateway VM setup

This Terraform module aims to provide a fully managed infrastructure that integrates Microsoft's Power Platform and Azure services. Utilizing both powerplatform and azurerm Terraform providers, this module encapsulates best practices and serves as a reference architecture for scalable, reliable, and manageable cloud infrastructure.

## Prerequisites

- Azure subscription
- SAP S/4HANA system
- Power Platform environment
- Service Principal with appropriate permissions (see below)

{{ .ModuleDetails }}

## Usage

Include this module in your Terraform scripts as follows:

```hcl

module "sap_gateway" {
  source            = "./quickstarts/{{ .ModuleName }}"
}

```

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `client_id_gw` | The client id / app id of the service principal where the on-premise data gateway admin permissions | string | `null` | true |
| `secret_gw` | The secret of the service principal with on-premise data gateway admin permissions | string | `null` | true |
| `tenant_id_gw` | The AAD tenant id of service principal or user | string | `null` | true |
| `subscription_id_gw` | The subscription id of the service principal with on-premise data gateway admin permissions | string | `null` | true |
| `region_gw` | The Azure region where the resources in this example should be created | string | `null` | true |
| `vm_pwd_gw` | The password for the VM of the on-premise data gateway | string | `null` | true |
| `prefix` | The prefix which should be used for all resources name | string | `opdgw` | false |
| `base_name` | The base name which should be used for all resources name | string | `AzureSAPIntegration` | false |
| `installps7_link` | The Blob link to the PowerShell 7 installation file | string | `null` | true |

## Output Values

## Detailed Behavior

### On-Premise Data Gateway

A Windows Virtual Machine is created with all the software connectors required to connect on-premise data gateway and self-hosted integration runtime.

It is the list of software installed on the Virtual Machine:

#### PowerShell 7

It is required to execute the script for the on-premise data gateway installation. After VM creation, the [script](./gateway-vm/scripts/installps7.ps1) download and install PowerShell 7.

#### Java Runtime

to-do: describe

#### Microsoft Data Connector

to-do: describe

#### SAP Connector for .Net

to-do: describe

### Gateway Principal

A service principal is created to allow the on-premise data gateway to connect the gateway back to Power Platform.  The service principal is granted the following permissions:

- `Application.Read.All`
- `Application.ReadWrite.All`

### Power Platform Resources

A PowerApps Environment is created

### Network

All resources are provisioned within the same Azure Virtual Network, ensuring that they can communicate securely without exposure to the public internet.

## Limitations and Considerations

- Due to Power Platform limitations, certain resources may not fully support Terraform's state management.
- Make sure to set appropriate RBAC for Azure and Power Platform resources.
- This module is provided as a sample only and is not intended for production use without further customization.

## Additional Resources

- [Power Platform Admin Documentation](https://learn.microsoft.com/en-us/power-platform/admin/)
- [Azure AD Terraform Provider](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_configuration)
- <https://learn.microsoft.com/en-us/power-platform/admin/wp-onpremises-gateway>
