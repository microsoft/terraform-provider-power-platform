output "coe_automation_satelite_solution_settings_empty_connection_reference_length" {
  value = local.empty_connection_reference_satelite_local
}

resource "powerplatform_environment" "satelite_environment" {
  display_name     = "AutimationKitEnvironmentSatelite"
  location         = "europe"
  language_name    = "1033"
  currency_name    = "USD"
  environment_type = "Sandbox"

  is_custom_controls_in_canvas_apps_enabled = true

  depends_on = [
    powerplatform_package.main_environment_deployment_package
  ]
}


resource "powerplatform_solution" "creator_kit_solution_satelite" {
  solution_file    = var.creator_kit_solution_zip_path
  settings_file    = ""
  solution_name    = "CreatorKitCore"
  environment_name = powerplatform_environment.satelite_environment.environment_name

  depends_on = [
  ]
}

## We only import the main solution when connection reference id are set in variables
resource "powerplatform_solution" "automation_kit_satelite_solution" {

  count = local.empty_connection_reference_satelite_local == 0 ? 1 : 0

  solution_file    = var.automation_coe_satellite_solution_zip_path
  settings_file    = "${path.module}${var.setting_file_path}"
  solution_name    = "AutomationCoESatellite"
  environment_name = powerplatform_environment.satelite_environment.environment_name

  depends_on = [
    powerplatform_solution.creator_kit_solution_satelite,
  ]
}

## We only import the main solution when connection reference id are set in variables
resource "powerplatform_package" "satelite_environment_deployment_package" {

  count = local.empty_connection_reference_satelite_local == 0 ? 1 : 0

  environment_name = powerplatform_environment.satelite_environment.environment_name
  package_name     = "Microsoft_AutomationKIT_Satellite_Package"
  package_file     = "${path.module}/data/Microsoft_AutomationKIT_Satellite_Package.zip"
  package_settings = "installsatellitesolution=False|importconfigdata=True|activateallflows=False"

  depends_on = [
    powerplatform_solution.automation_kit_satelite_solution,
    powerplatform_user.power_platform_automation_kit_app_user
  ]
}

#assign app user to the satelite env
resource "powerplatform_user" "power_platform_automation_kit_app_user" {
  environment_name = powerplatform_environment.satelite_environment.environment_name
  //aad_id = var.automation_kit_application_id
  is_app_user    = true
  application_id = var.automation_kit_application_id
  first_name     = "#" #dataverse will force # as the app user name
  last_name      = "CoE_Automation_Kit_App_User"
  security_roles = ["System Administrator"]

  depends_on = [
  ]
}
#####################

#find app id that was deployed to main env
data "powerplatform_powerapps" "environment_app_filter" {
  environment_name = powerplatform_environment.main_environment.environment_name

  depends_on = [
    powerplatform_solution.automation_kit_main_solution,
    powerplatform_package.main_environment_deployment_package
  ]
}

locals {
  only_specific_name = toset(
    [
      for each in data.powerplatform_powerapps.environment_app_filter.apps :
      each if each.display_name == "Automation Project"
  ])
}
#####################



locals {
  all_connection_reference_satelite_local = "${var.satelite_conn_ref_shared_office365}${var.satelite_conn_ref_shared_powerplatformforadmins}${var.satelite_conn_ref_shared_commondataserviceforapps}${var.satelite_conn_ref_shared_commondataservice}${var.satelite_conn_ref_shared_flowmanagement}${var.satelite_conn_ref_shared_office365users}"
  empty_connection_reference_satelite_local = (length(local.all_connection_reference_satelite_local) -
  length(replace(local.all_connection_reference_satelite_local, "replace_with_your_connection_reference", "")))
}

variable "setting_file_path" {
  default = "/data/coe_automation_satelite_solution_settings.json"
}

resource "local_file" "coe_automation_satelite_solution_settings" {
  filename = "${path.module}${var.setting_file_path}"
  content  = <<EOF
{
  "EnvironmentVariables": [
    {
      "SchemaName": "autocoe_AKVClientIdSecret",
      "Value": "${var.key_vault_name}/secrets/${var.key_vault_secret_client_id_name}"
    },
    {
      "SchemaName": "autocoe_AKVClientSecretSecret",
      "Value": "${var.key_vault_name}/secrets/${var.key_vault_secret_client_password_name}"
    },
    {
      "SchemaName": "autocoe_AKVTenantIdSecret",
      "Value": "${var.key_vault_name}/secrets/${var.key_vault_client_secret_tenant_id_name}"
    },
    {
      "SchemaName": "autocoe_AutomationCoEAlertEmailRecipient",
      "Value": "${var.env_variable_autocoe_AutomationCoEAlertEmailRecipient}"
    },
    {
      "SchemaName": "autocoe_AutomationProjectAppID",
      "Value": "${length(local.only_specific_name) > 0 ? one(local.only_specific_name).name : ""}"
    },
    {
      "SchemaName": "autocoe_DesktopFlowsBaseURL",
      "Value": "https://make.powerautomate.com/environments/"
    },
    {
      "SchemaName": "autocoe_EnvironmentId",
      "Value": "${powerplatform_environment.satelite_environment.environment_name}"
    },
    {
      "SchemaName": "autocoe_EnvironmentName",
      "Value": "${powerplatform_environment.satelite_environment.display_name}"
    },
    {
      "SchemaName": "autocoe_EnvironmentRegion",
      "Value": "${powerplatform_environment.satelite_environment.location}"
    },
    {
      "SchemaName": "autocoe_EnvironmentUniqeName",
      "Value":  "${powerplatform_environment.satelite_environment.environment_name}"
    },
    {
      "SchemaName": "autocoe_EnvironmentUniqueNameofCoEMain",
      "Value":  "${powerplatform_environment.main_environment.environment_name}"
    },
    {
      "SchemaName": "autocoe_EnvironmentURL",
      "Value": "https://${powerplatform_environment.satelite_environment.url}.${powerplatform_environment.satelite_environment.domain}.dynamics.com/"
    },
    {
      "SchemaName": "autocoe_FlowSessionTraceRecordOwnerId",
      "Value": "${var.env_variable_autocoe_FlowSessionTraceRecordOwnerId}"
    },
    {
      "SchemaName": "autocoe_StoreExtractedScript",
      "Value": "${var.env_variable_autocoe_StoreExtractedScript}"
    }
  ],
  "ConnectionReferences": [
    {
      "LogicalName": "autocoe_DataverseAutoCoESatellitecurrent",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_DataverseLegacyAutoCoEMain",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataservice}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
    },
    {
      "LogicalName": "autocoe_DLPImpactAnalysisDataverse",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_sharedcommondataserviceforapps_51023",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_sharedcommondataserviceforapps_6489c",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_sharedcommondataserviceforapps_98ee0",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_sharedcommondataservice_305b1",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataservice}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
    },
    {
      "LogicalName": "autocoe_sharedcommondataservice_34c59",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataservice}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
    },
    {
      "LogicalName": "autocoe_sharedcommondataservice_42f40",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataservice}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
    },
    {
      "LogicalName": "autocoe_sharedcommondataservice_9a93f",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataservice}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
    },
    {
      "LogicalName": "autocoe_sharedcommondataservice_a6d4e",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataservice}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
    },
    {
      "LogicalName": "autocoe_sharedflowmanagement_2d016",
      "ConnectionId": "${var.satelite_conn_ref_shared_flowmanagement}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_flowmanagement"
    },
    {
      "LogicalName": "autocoe_sharedflowmanagement_85ee3",
      "ConnectionId": "${var.satelite_conn_ref_shared_flowmanagement}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_flowmanagement"
    },
    {
      "LogicalName": "autocoe_sharedflowmanagement_ccea3",
      "ConnectionId": "${var.satelite_conn_ref_shared_flowmanagement}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_flowmanagement"
    },
    {
      "LogicalName": "autocoe_sharedflowmanagement_e5dfb",
      "ConnectionId": "${var.satelite_conn_ref_shared_flowmanagement}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_flowmanagement"
    },
    {
      "LogicalName": "autocoe_sharedoffice365users_7aa15",
      "ConnectionId": "${var.satelite_conn_ref_shared_office365users}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_office365users"
    },
    {
      "LogicalName": "autocoe_sharedpowerplatformforadmins_149dd",
      "ConnectionId": "${var.satelite_conn_ref_shared_powerplatformforadmins}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_powerplatformforadmins"
    },
    {
      "LogicalName": "cr50e_sharedoffice365_abd1a",
      "ConnectionId": "${var.satelite_conn_ref_shared_office365}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_office365"
    },
    {
      "LogicalName": "new_sharedcommondataservice_38df0",
      "ConnectionId": "${var.satelite_conn_ref_shared_commondataservice}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
    },
    {
      "LogicalName": "new_sharedoffice365_912a8",
      "ConnectionId": "${var.satelite_conn_ref_shared_office365}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_office365"
    }
  ]
}
EOF
}
