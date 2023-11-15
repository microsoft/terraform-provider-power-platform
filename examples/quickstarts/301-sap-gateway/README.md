<!-- This document is auto-generated. Do not edit directly. Make changes to README.md.tmpl instead. -->
# SAP Connectivity Runtime Setup (301 level)

This Terraform module aims to provide a fully managed infrastructure that integrates Microsoft's Power Platform and Azure services with SAP Systems. Utilizing  `azurerm` and `azurecaf` Terraform providers, this module encapsulates best practices and serves as a reference architecture for scalable, reliable, and manageable cloud infrastructure.

## Prerequisites

- Service Principal or User Account with permissions configured as referenced in [the provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication).

### Storage Account Preparation

Before you execute the script, you need to upload the SAP .NET Connector MSI file the folder `./storage-account/sapnco-msi` and rename to `sapnco.msi` (check below for more information).

### SHIR Nodes Preparation

Make sure there is not any node assigned to the self-hosted integration runtime at Synapse or ADF.

## Example Files

The example files can be found in `examples/quickstarts/301-sap-gateway`

## Provider Requirements:
* **azurecaf (`aztfmod/azurecaf`):** `>=1.2.26`
* **azurerm (`hashicorp/azurerm`):** `>=3.74.0`
* **random:** (any version)

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `base_name` | The base name which should be used for all resources name | string | `"AzureSAPIntegration"` | false |
| `client_id_gw` | The client ID / app ID of the service principal where the on-premise data gateway admin permissions | string | `null` | true |
| `client_id_pp` | The client ID / app ID of the service principal with Power Platform admin permissions | string | `null` | true |
| `gateway_name` | The name of the gateway to be created on Power Platform | string | `null` | true |
| `prefix` | The prefix which should be used for all resources name | string | `"opdgw"` | false |
| `recover_key_gw` | The recovery key of the gateway | string | `null` | true |
| `region_gw` | The Azure region where the resources in this example should be created | string | `null` | true |
| `sap_subnet_id` | The SAP system subnet ID | string | `null` | true |
| `secret_gw` | The secret of the service principal with on-premise data gateway admin permissions | string | `null` | true |
| `secret_pp` | The secret of the service principal with Power Platform admin permissions | string | `null` | true |
| `shir_key` | Value of the secret name for the IR key | string | `null` | true |
| `subscription_id_gw` | The subscription ID of the service principal with on-premise data gateway admin permissions | string | `null` | true |
| `tenant_id_gw` | The tenant ID of service principal or user | string | `null` | true |
| `tenant_id_pp` | The tenant ID of service principal or user at Power Platform | string | `null` | true |
| `user_id_admin_pp` | The user ID to be assigned as Admin role of the Power Platform | string | `null` | true |


## Resources
* `azurecaf_name.key_vault` from `azurecaf`
* `azurecaf_name.key_vault_secret_irkey` from `azurecaf`
* `azurecaf_name.key_vault_secret_pp` from `azurecaf`
* `azurecaf_name.key_vault_secret_recover_key` from `azurecaf`
* `azurecaf_name.key_vault_secret_vm_pwd` from `azurecaf`
* `azurecaf_name.nic` from `azurecaf`
* `azurecaf_name.nsg` from `azurecaf`
* `azurecaf_name.publicip` from `azurecaf`
* `azurecaf_name.rg` from `azurecaf`
* `azurecaf_name.subnet` from `azurecaf`
* `azurecaf_name.vnet` from `azurecaf`
* `azurerm_key_vault.key_vault` from `azurerm`
* `azurerm_key_vault_access_policy.key_vault_access_policy` from `azurerm`
* `azurerm_key_vault_secret.key_vault_secret_irkey` from `azurerm`
* `azurerm_key_vault_secret.key_vault_secret_pp` from `azurerm`
* `azurerm_key_vault_secret.key_vault_secret_recover_key` from `azurerm`
* `azurerm_key_vault_secret.key_vault_secret_vm_pwd` from `azurerm`
* `azurerm_network_interface.nic` from `azurerm`
* `azurerm_network_interface_security_group_association.rgassociation` from `azurerm`
* `azurerm_network_security_group.nsg` from `azurerm`
* `azurerm_public_ip.publicip` from `azurerm`
* `azurerm_resource_group.rg` from `azurerm`
* `azurerm_subnet.subnet` from `azurerm`
* `azurerm_virtual_network.vnet` from `azurerm`
* `random_string.key_vault_suffix` from `random`
* `random_string.vm_pwd` from `random`

## Data Sources
* `data.azurerm_client_config.current` from `azurerm`

## Child Modules
* `gateway_vm` from `./gateway-vm`
* `storage_account` from `./storage-account`


## Usage

The entire script is required for the proper installation, unless you decide to create any one of the resources separatelly.

You have to execute the normal Terraform commands:

``terraform init -upgrade

``terraform plan -var-file=local.tfvars

``terraform apply -var-file=local.tfvars

``terraform destroy -var-file=local.tfvars

## Terraform Version Constraints

- azurerm `>=3.74.0`
- azurecaf `>=1.2.26`

## Detailed Behavior

### On-Premise Data Gateway

A Windows Virtual Machine is created with all the software connectors required to connect on-premise data gateway and self-hosted integration runtime.

It is the list of software installed on the Virtual Machine:

#### PowerShell 7

It is required to execute the script for the on-premise data gateway installation. After VM creation, the [script](./gateway-vm/scripts/ps7-setup.ps1) download and install PowerShell 7.

#### Java Runtime

It is required for SHIR runtime and SAP data flows. Check the [prerequisites](https://learn.microsoft.com/en-us/azure/data-factory/create-self-hosted-integration-runtime?tabs=data-factory#prerequisites) in the documentation for more details.

#### Microsoft Self-Hosted Integration Runtime (SHIR)

It is the runtime used to connect the VM to SHIR in Synapse/ADF/Fabric. Check the [documentation](https://learn.microsoft.com/en-us/azure/data-factory/create-self-hosted-integration-runtime) for more details.

#### SAP Connector for .Net

It is the runtime used to connect the VM to SAP system. You need to [download the MSI file](https://support.sap.com/en/product/connectors/msnet.html) and upload to the folder mentioned above. The installation provided in this guide, follows this [documentation](https://learn.microsoft.com/en-us/azure/data-factory/sap-change-data-capture-shir-preparation).

#### On-Premises Data Gateway

It is the runtime used to connect to the Power Platform connectors (e.g. SAP ERP). Here is some references used to created the script:

- [Learn how to install On-premises data gateway for Azure Analysis Services | Microsoft Learn](https://learn.microsoft.com/en-us/azure/analysis-services/analysis-services-gateway-install?tabs=azure-powershell)
- [Data Gateway Documentation](https://learn.microsoft.com/en-us/powershell/module/datagateway/?view=datagateway-ps)

### Power Platform Resources

A PowerApps Environment is created.

### Network

All resources are provisioned within the same Azure Virtual Network where the SAP System is installed (`sap_subnet_id` input parameter), ensuring that they can communicate securely without exposure to the public internet.

## Limitations and Considerations

- Due to Power Platform limitations, certain resources may not fully support Terraform's state management.
- Make sure to set appropriate RBAC for Azure and Power Platform resources.
- This module is provided as a sample only and is not intended for production use without further customization.

## Additional Resources

- [Power Platform Admin Documentation](https://learn.microsoft.com/en-us/power-platform/admin/)
- [Azure AD Terraform Provider](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_configuration)
- <https://learn.microsoft.com/en-us/power-platform/admin/wp-onpremises-gateway>
