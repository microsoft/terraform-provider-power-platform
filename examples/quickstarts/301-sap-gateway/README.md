<!-- This document is auto-generated. Do not edit directly. Make changes to README.md.tmpl instead. -->
# SAP Connectivity Runtime Setup (301 level)

This Terraform module aims to provide a fully managed infrastructure that integrates Microsoft's Power Platform and Azure services with SAP Systems. Utilizing  `azurerm` and `azurecaf` Terraform providers, this module encapsulates best practices and serves as a reference architecture for scalable, reliable, and manageable cloud infrastructure.

To provide connectivity between SAP and Microsoft services, it is required to install a set of applications and also setup proper configuration. This document is a guide to setup a Terraform script provisioning the Virtual Machine and all the requirements to connect the SHIR (Self-hosted Integration Runtime), the Microsoft Gateway and SAP .NET Connector.

## Prerequisites and Preparation

### Credentials

- Azure subscription
- Service Principal or User Account with permissions configured as referenced in [the provider's user documentation](https://microsoft.github.io/terraform-provider-power-platform#authentication).
- All the credentials for Azure resources creation.

### SAP Systems

For the execution of this Terraform script, you do not need the SAP credentials or the application server information. However, it is important to know how to connect to the SAP system to provide proper information to the script.

Check more details in the section about networking below.

### MSI file for SAP .NET Connector

Before you execute the script, you need to copy the SAP .NET Connector MSI file the folder `./storage-account/sapnco-msi` and rename to `sapnco.msi`.

The SAP Connector for .NET is the runtime used to connect the VM to SAP system. You need to [download the MSI file](https://support.sap.com/en/product/connectors/msnet.html) and upload to the folder mentioned above. The installation provided in this guide, follows this [documentation](https://learn.microsoft.com/en-us/azure/data-factory/sap-change-data-capture-shir-preparation).

> Note: The MSI file cannot be downloaded automatically, because it is required SAP S-User authentication.

### SHIR Nodes Preparation

Make sure there is not any node assigned to the self-hosted integration runtime at Synapse or ADF. Please check more information at [Create a self-hosted integration runtime - Azure Data Factory & Azure Synapse | Microsoft Learn](https://learn.microsoft.com/en-us/azure/data-factory/create-self-hosted-integration-runtime?tabs=data-factory).

### Networking Requirements

The Terraform code on this repo will deploy the VM on an existing Azure Vnet so you need to provide the Azure "sap_subnet_id" value, check more details in the section about networking below.

## Terraform Version Constraints

- azurerm `>=3.74.0`
- azurecaf `>=1.2.26`

## Example Files

The example files can be found in `examples/quickstarts/301-sap-gateway`

## Provider Requirements

- **azurecaf (`aztfmod/azurecaf`):** `>=1.2.26`
- **azurerm (`hashicorp/azurerm`):** `>=3.74.0`
- **random:** (any version)

## Input Variables

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| `base_name` | The base name which should be used for all resources name | string | `"AzureSAPIntegration"` | false |
| `client_id_gw` | The client ID / app ID of the service principal where the on-premise data gateway admin permissions | string | `null` | true |
| `client_id_pp` | The client ID / app ID of the service principal with Power Platform admin permissions | string | `null` | true |
| `gateway_name` | The name of the gateway to be created on Power Platform | string | `null` | true |
| `ir_key` | Value of the secret name for the IR key | string | `null` | true |
| `prefix` | The prefix which should be used for all resources name | string | `"opdgw"` | false |
| `recover_key_gw` | The recovery key of the gateway | string | `null` | true |
| `region_gw` | The Azure region where the resources in this example should be created | string | `null` | true |
| `sap_subnet_id` | The SAP system subnet ID | string | `null` | true |
| `secret_gw` | The secret of the service principal with on-premise data gateway admin permissions | string | `null` | true |
| `secret_pp` | The secret of the service principal with Power Platform admin permissions | string | `null` | true |
| `subscription_id_gw` | The subscription ID of the service principal with on-premise data gateway admin permissions | string | `null` | true |
| `tenant_id_gw` | The tenant ID of service principal or user | string | `null` | true |
| `tenant_id_pp` | The tenant ID of service principal or user at Power Platform | string | `null` | true |
| `user_id_admin_pp` | The user ID to be assigned as Admin role of the Power Platform | string | `null` | true |

## Resources

- `azurecaf_name.key_vault` from `azurecaf`
- `azurecaf_name.key_vault_secret_irkey` from `azurecaf`
- `azurecaf_name.key_vault_secret_pp` from `azurecaf`
- `azurecaf_name.key_vault_secret_recover_key` from `azurecaf`
- `azurecaf_name.key_vault_secret_vm_pwd` from `azurecaf`
- `azurecaf_name.nic` from `azurecaf`
- `azurecaf_name.nsg` from `azurecaf`
- `azurecaf_name.publicip` from `azurecaf`
- `azurecaf_name.rg` from `azurecaf`
- `azurecaf_name.subnet` from `azurecaf`
- `azurecaf_name.vnet` from `azurecaf`
- `azurerm_key_vault.key_vault` from `azurerm`
- `azurerm_key_vault_access_policy.key_vault_access_policy` from `azurerm`
- `azurerm_key_vault_secret.key_vault_secret_irkey` from `azurerm`
- `azurerm_key_vault_secret.key_vault_secret_pp` from `azurerm`
- `azurerm_key_vault_secret.key_vault_secret_recover_key` from `azurerm`
- `azurerm_key_vault_secret.key_vault_secret_vm_pwd` from `azurerm`
- `azurerm_network_interface.nic` from `azurerm`
- `azurerm_network_interface_security_group_association.rgassociation` from `azurerm`
- `azurerm_network_security_group.nsg` from `azurerm`
- `azurerm_public_ip.publicip` from `azurerm`
- `azurerm_resource_group.rg` from `azurerm`
- `azurerm_subnet.subnet` from `azurerm`
- `azurerm_virtual_network.vnet` from `azurerm`
- `random_string.key_vault_suffix` from `random`
- `random_string.vm_pwd` from `random`

## Data Sources

- `data.azurerm_client_config.current` from `azurerm`

## Child Modules

- `gateway_vm` from `./gateway-vm`
- `storage_account` from `./storage-account`

## Resources to Be Created

Here is a diagram of all resources to be created:

![Components to be installed](./.img/components.svg)

### Resource Group

All the resources are created in the same resource group.

### Network

All resources are provisioned within the same Azure Virtual Network and it is connected to the SAP system vnet using the variable `sap_subnet_id`, ensuring that they can communicate securely without exposure to the public internet.

The scenario considered here is the SAP system installed on-premisses on Azure, and you cannot use a public address. So, you will need to provide the subnet ID where the SAP system is installed.

The subnet ID is available at the JSON view of the SAP system virtual network, in the parameter `id`. It is expected something like below:

`/subscriptions/abababab-12ab-ab00-82e2-aa00babab102/resourceGroups/resource-group-name/providers/Microsoft.Network/virtualNetworks/VNet-name/subnets/default`

### Key Vault

The Key Vault contains the credentials required to connect all the components.

| Secret Name               | Description                                                                      |
| ------------------------- | -------------------------------------------------------------------------------- |
| \<prefix>-kvs-pp-xxx         | Service Principal secrete for Power Platforms where the Gateway will be created. |
| \<prefix>-kvs-irkey-xxx      | IR Key                                                                           |
| \<prefix>-kvs-recoverkey-xxx | Recovery key used during SHIR configuration                                      |
| \<prefix>-kvs-vm-pwd-xxx     | VM password (randomly generated)                                                 |

> Notes: "xxx" corresponding to the random generated part of the name.
> You need to add an access policy with proper key permissions for your user to access the values created by the script.

### Storage Account

The module `storage-account` creates the storage account and upload the installation files and installation scripts. All the scripts used to install the required applications, are in the folder `./storage-account/scripts`.

### Virtual Machine

The module `gateway-vm` creates a Windows Server 2022 virtual machine and install the following applications listed below. You can use any other Windows version following the requirements for each one of the applications listed below.

> Note: To access the VM using Remote Desktop you can use the user `sapadmin` and the password saved in the secret at the Key Vault.

Here is the detail for every application installed in the virtual machine.

#### PowerShell 7

The module `ps7-setup` installs the PowerShell version 7. It is required to execute the script for the on-premise data gateway installation. After VM creation, the script [ps7-setup.ps1](./storage-account/scripts/ps7-setup.ps1) download and install the PowerShell 7.

This module creates an Azure Application Definition which upload the installation script to the VM and the Azure Application Version executes that script.

#### Java Runtime

The Parquet file format is required for SHIR runtime and SAP data flows execute delta extractions. So allow the SHIR use e Parquet file format, it is required to install the Java Runtime. Check the [prerequisites](https://learn.microsoft.com/en-us/azure/data-factory/create-self-hosted-integration-runtime?tabs=data-factory#prerequisites) in the documentation for more details.

The module `java-runtime-setup` installs the Java Runtime. After VM creation, the script [java-setup.ps1](./storage-account/scripts/java-setup.ps1) download and install the Java Runtime.

This module creates an Azure Application Definition which upload the installation script to the VM and the Azure Application Version executes that script.

#### Microsoft Self-Hosted Integration Runtime (SHIR)

It is the runtime used to connect the VM to SHIR in Synapse/ADF/Fabric. Check the [documentation](https://learn.microsoft.com/en-us/azure/data-factory/create-self-hosted-integration-runtime) for more details.

The module `runtime-setup` installs the SHIR runtime and connect to the Integration Runtime using its key code. You should specify the key using the variable `ir_key` in the file `local.tfvars`.

#### SAP Connector for .Net

It is the runtime used to connect the VM to SAP system. You need to [download the MSI file](https://support.sap.com/en/product/connectors/msnet.html) and upload to the folder mentioned above. The installation provided in this guide, follows this [documentation](https://learn.microsoft.com/en-us/azure/data-factory/sap-change-data-capture-shir-preparation).

#### On-Premises Data Gateway

The module `runtime-setup` installs the the On-Premises Data Gateway. It is required to connect to the Power Platform connectors (e.g. SAP ERP) to the SAP systems.

The configuration script reads the Power Platforms secret and recovery key from the Key Vault, and creates the gateway in the PowerApps website using the name provided in the variable `gateway_name` in the file `local.tfvars`. The script also needs the the variable `user_id_admin_pp` to add user ID as the gateway admin. So, that user will be able to access the Gateway in the PowerApps website.

Here is some references used to created the script:

- [Learn how to install On-premises data gateway for Azure Analysis Services | Microsoft Learn](https://learn.microsoft.com/en-us/azure/analysis-services/analysis-services-gateway-install?tabs=azure-powershell)
- [Data Gateway Documentation](https://learn.microsoft.com/en-us/powershell/module/datagateway/?view=datagateway-ps)

### Naming Convention

The majority of the resources names are following a name convention for clear identification and ensure unique name. The name convention is composed by the prefix, resource type name, base name and random 3-char text. You can change the prefix and base name in the file `local.tfvars` using respective variables.

## Limitations and Considerations

- Due to Power Platform limitations, certain resources may not fully support Terraform's state management.
- Make sure to set appropriate RBAC for Azure and Power Platform resources.
- This module is provided as a sample only and is not intended for production use without further customization.

## Usage

The entire script is required for the proper installation, unless you decide to create any one of the resources separately.

You have to execute the normal Terraform commands:

```bash
terraform init -upgrade
```

The command `terraform init` is used to initialize a working directory that contains Terraform configuration files. It performs several tasks, such as configuring the backend for storing the state, installing the required providers and modules, and creating a lock file to track the versions of the providers and modules. This command is the first step in the Terraform workflow and should be run whenever the configuration changes or a new workspace is created. It is safe to run this command multiple times, as it will not delete or overwrite any existing configuration or state.

```bash
terraform plan -var-file=local.tfvars
```

This command tells Terraform to create an execution plan that shows the changes that Terraform will make to the infrastructure based on the configuration files in the current working directory and the variable values in the file `local.tfvars`. The -var-file option allows the user to specify a file that contains variable definitions for the root module of the configuration. You can find an example of this file in the next section below.

The execution plan consists of a set of actions that will create, update, or destroy resources. The plan allows the user to review the changes before applying them, or to save the plan to a file for later use. The plan also helps to ensure that the configuration and the state are in sync, and that the changes match the user’s expectations

```bash
terraform apply -var-file=local.tfvars
```

The `terraform apply` command creates or updates infrastructure depending on the configuration files. By default, a plan will be generated first and will need interactive approval before applying. The plan shows the actions that Terraform will take to create, modify, or destroy resources.

```bash
terraform destroy -var-file=local.tfvars
```

The terraform destroy command is a convenient way to destroy all remote objects managed by a particular Terraform configuration. It is the opposite of the terraform apply command. The terraform destroy command is useful when Terraform is used to manage ephemeral infrastructure for development purposes, and the user wants to clean up all of the temporary objects once the work is finished.

### Example of local.tfvars file

Here is an example of the `local.tfvars` file.

```bash
client_id_gw       = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
secret_gw          = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
tenant_id_gw       = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
subscription_id_gw = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
region_gw          = "West Europe"
client_id_pp       = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
secret_pp          = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
tenant_id_pp       = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
user_id_admin_pp   = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
ir_key             = "IR@XXXXXX-XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
gateway_name       = "OPDGW-SAPAzureIntegration"
recover_key_gw     = "XXXXXXXXXXX"
sap_subnet_id      = "/subscriptions/abababab-12ab-ab00-82e2-aa00babab102/resourceGroups/resource-group-name/providers/Microsoft.Network/virtualNetworks/VNet-name/subnets/default"
prefix             = "opdgw"
base_name          = "AzureSAPIntegration"
```

> Note: You need to provide all the "local.tfvars" values to run the command "terraform apply".

## Additional Resources

- [Overview of VM Applications in the Azure Compute Gallery - Azure Virtual Machines | Microsoft Learn](https://learn.microsoft.com/en-us/azure/virtual-machines/vm-applications)
- [Create and configure a self-hosted integration runtime](https://learn.microsoft.com/en-us/azure/data-factory/create-self-hosted-integration-runtime)
- [Power Platform Admin Documentation](https://learn.microsoft.com/en-us/power-platform/admin/)
- [Azure AD Terraform Provider](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_configuration)
- <https://learn.microsoft.com/en-us/power-platform/admin/wp-onpremises-gateway>
