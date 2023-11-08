# SAP Connectivity Runtime Setup

This Terraform module aims to provide a fully managed infrastructure that integrates Microsoft's Power Platform and Azure services with SAP Systems. Utilizing  `azurerm` and `azurecaf` Terraform providers, this module encapsulates best practices and serves as a reference architecture for scalable, reliable, and manageable cloud infrastructure.

## Prerequisites

- Service Principal or User Account with permissions configured as referenced in [the provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication).
- [SAP S/4HANA system](https://cal.sap.com/)

### Storage Account

Before you execute the script, you need to create a Storage Account, upload all scripts listed in the folder `./gateway-vm/scripts` and the SAP .NET Connector MSI file (check below for more information).

You can follow [this guide](https://learn.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal) to create the storage account.

> Tip: Make sure the blob accept anonymous access.

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

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `client_id_gw` | The client id / app id of the service principal where the on-premise data gateway admin permissions | string | `null` | true |
| `secret_gw` | The secret of the service principal with on-premise data gateway admin permissions | string | `null` | true |
| `tenant_id_gw` | The tenant id of service principal or user | string | `null` | true |
| `subscription_id_gw` | The subscription id of the service principal with on-premise data gateway admin permissions | string | `null` | true |
| `region_gw` | The Azure region where the resources in this example should be created | string | `null` | true |
| `sap_subnet_id` | The SAP system subnet ID | string | `null` | true |
| `vm_pwd_gw` | The password for the VM of the on-premise data gateway | string | `null` | true |
| `client_id_pp` | The client id / app id of the service principal with Power Platform admin permissions | string | `null` | true |
| `secret_pp` | The secret of the service principal with Power Platform admin permissions | string | `null` | true |
| `tenant_id_pp` | The tenant id of service principal or user at Power Platform | string | `null` | true |
| `ps7_setup_link` | The Blob link to the PowerShell 7 installation file | string | `null` | true |
| `java_setup_link` | The Blob link to the Java Runtime installation file | string | `null` | true |
| `user_id_admin_pp` | The user id to be assigned as Admin role of the Power Platform | string | `null` | true |
| `sapnco_install_link` | The Blob link to the SAP NCo installation file | string | `null` | true |
| `runtime_setup_link` | The Blob link to the runtime setup script | string | `null` | true |
| `shir_key` | Value of the secret name for the IR key | string | `null` | true |
| `gateway_name` | The name of the gateway to be created on Power Platform | string | `null` | true |
| `recover_key_gw` | The recovery key of the gateway | string | `null` | true |
| `prefix` | The prefix which should be used for all resources name | string | `opdgw` | false |
| `base_name` | The base name which should be used for all resources name | string | `AzureSAPIntegration` | false |

## Output Values

No output.

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

It is the runtime used to connect the VM to SAP system. You need to [download the MSI file](https://support.sap.com/en/product/connectors/msnet.html) and upload to the Storage Account mentioned above. The installation provided in this guide, follows this [documentation](https://learn.microsoft.com/en-us/azure/data-factory/sap-change-data-capture-shir-preparation).

#### On-Premises Data Gateway

It is the runtime used to connect to the Power Platform connectors (e.g. SAP ERP). Here is some references used to created the script:

- [Learn how to install On-premises data gateway for Azure Analysis Services | Microsoft Learn](https://learn.microsoft.com/en-us/azure/analysis-services/analysis-services-gateway-install?tabs=azure-powershell)
- [Data Gateway Documentation](https://learn.microsoft.com/en-us/powershell/module/datagateway/?view=datagateway-ps)

### Gateway Principal (future releases)

A service principal is created to allow the on-premise data gateway to connect the gateway back to Power Platform.  The service principal is granted the following permissions:

- `Application.Read.All`
- `Application.ReadWrite.All`

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
