output "coe_automation_main_solution_settings_empty_connection_reference_length" {
  value = local.empty_connection_reference_main_local
}

resource "powerplatform_environment" "main_environment" {

  display_name = "AutimationKitEnvironmentMain"
  location = "europe"
  language_name = "1033"
  currency_name = "USD"
  environment_type = "Sandbox"

  is_custom_controls_in_canvas_apps_enabled = true
}

resource "powerplatform_solution" "creator_kit_solution_main" {
  solution_file = var.creator_kit_solution_zip_path
  settings_file = ""
  solution_name = "CreatorKitCore"
  environment_name = powerplatform_environment.main_environment.environment_name
}

## We only import the main solution when connection reference id are set in variables
resource "powerplatform_solution" "automation_kit_main_solution" {

  count = local.empty_connection_reference_main_local == 0 ? 1 : 0
  
  solution_file = var.automation_coe_main_solution_zip_path
  settings_file = local_file.coe_automation_main_solution_settings.filename
  solution_name = "AutomationCoEMain"
  environment_name = powerplatform_environment.main_environment.environment_name

  depends_on = [
    powerplatform_solution.creator_kit_solution_main
  ]
}

## We only import the main solution when connection reference id are set in variables
resource "powerplatform_package" "main_environment_deployment_package" {

  count = local.empty_connection_reference_main_local == 0 ? 1 : 0
  
  environment_name = powerplatform_environment.main_environment.environment_name
  package_name = "Microsoft_AutomationKIT_Main_Package"
  package_file = "${path.module}/data/Microsoft_AutomationKIT_Main_Package.zip"
  //package_settings = "installmainsolution=true|importconfigdata=True|AutomationCoEMain_componentarguments=W3siQG9kYXRhLnR5cGUiOiJNaWNyb3NvZnQuRHluYW1pY3MuQ1JNLmNvbm5lY3Rpb25yZWZlcmVuY2UiLCJjb25uZWN0aW9uUmVmZXJlbmNlTG9naWNhbE5hbWUiOiJhdXRvY29lX3NoYXJlZGFwcHJvdmFsc184ODY1ZCIsImNvbm5lY3Rpb25SZWZlcmVuY2VEaXNwbGF5TmFtZSI6IkFwcHJvdmFscyBBdXRvbWF0aW9uQ29FQ29yZU1haW5Qcml2YXRlUHJldmlldy04ODY1ZCIsImRlc2NyaXB0aW9uIjpudWxsLCJjb25uZWN0b3JJZCI6Ii9wcm92aWRlcnMvTWljcm9zb2Z0LlBvd2VyQXBwcy9hcGlzL3NoYXJlZF9hcHByb3ZhbHMiLCJjb25uZWN0aW9uSWQiOiIwZjc0NWMyYzIxMjE0Mjc2YjRmMDNjOGVmMzdjMWMxYiJ9LHsiQG9kYXRhLnR5cGUiOiJNaWNyb3NvZnQuRHluYW1pY3MuQ1JNLmNvbm5lY3Rpb25yZWZlcmVuY2UiLCJjb25uZWN0aW9uUmVmZXJlbmNlTG9naWNhbE5hbWUiOiJhdXRvY29lX3NoYXJlZGNvbW1vbmRhdGFzZXJ2aWNlZm9yYXBwc185YzYzNyIsImNvbm5lY3Rpb25SZWZlcmVuY2VEaXNwbGF5TmFtZSI6Ik1pY3Jvc29mdCBEYXRhdmVyc2UgQXV0b21hdGlvbkNvRUNvcmVNYWluUHJpdmF0ZVByZXZpZXctOWM2MzciLCJkZXNjcmlwdGlvbiI6bnVsbCwiY29ubmVjdG9ySWQiOiIvcHJvdmlkZXJzL01pY3Jvc29mdC5Qb3dlckFwcHMvYXBpcy9zaGFyZWRfY29tbW9uZGF0YXNlcnZpY2Vmb3JhcHBzIiwiY29ubmVjdGlvbklkIjoiOTcwMTExZWM1MDMzNDI5ZjgxYjA2MjEzZTNjZGVkN2UifSx7IkBvZGF0YS50eXBlIjoiTWljcm9zb2Z0LkR5bmFtaWNzLkNSTS5jb25uZWN0aW9ucmVmZXJlbmNlIiwiY29ubmVjdGlvblJlZmVyZW5jZUxvZ2ljYWxOYW1lIjoiYXV0b2NvZV9zaGFyZWRjb21tb25kYXRhc2VydmljZWZvcmFwcHNfYjM2OTgiLCJjb25uZWN0aW9uUmVmZXJlbmNlRGlzcGxheU5hbWUiOiJNaWNyb3NvZnQgRGF0YXZlcnNlIEF1dG9tYXRpb25Db0VDb3JlTWFpblByaXZhdGVQcmV2aWV3LWIzNjk4IiwiZGVzY3JpcHRpb24iOm51bGwsImNvbm5lY3RvcklkIjoiL3Byb3ZpZGVycy9NaWNyb3NvZnQuUG93ZXJBcHBzL2FwaXMvc2hhcmVkX2NvbW1vbmRhdGFzZXJ2aWNlZm9yYXBwcyIsImNvbm5lY3Rpb25JZCI6Ijk3MDExMWVjNTAzMzQyOWY4MWIwNjIxM2UzY2RlZDdlIn0seyJAb2RhdGEudHlwZSI6Ik1pY3Jvc29mdC5EeW5hbWljcy5DUk0uY29ubmVjdGlvbnJlZmVyZW5jZSIsImNvbm5lY3Rpb25SZWZlcmVuY2VMb2dpY2FsTmFtZSI6ImF1dG9jb2Vfc2hhcmVkY29tbW9uZGF0YXNlcnZpY2Vmb3JhcHBzX2I3MWE0IiwiY29ubmVjdGlvblJlZmVyZW5jZURpc3BsYXlOYW1lIjoiTWljcm9zb2Z0IERhdGF2ZXJzZSBVcGRhdGVNYWNoaW5lU3RhdHVzIiwiZGVzY3JpcHRpb24iOm51bGwsImNvbm5lY3RvcklkIjoiL3Byb3ZpZGVycy9NaWNyb3NvZnQuUG93ZXJBcHBzL2FwaXMvc2hhcmVkX2NvbW1vbmRhdGFzZXJ2aWNlZm9yYXBwcyIsImNvbm5lY3Rpb25JZCI6Ijk3MDExMWVjNTAzMzQyOWY4MWIwNjIxM2UzY2RlZDdlIn0seyJAb2RhdGEudHlwZSI6Ik1pY3Jvc29mdC5EeW5hbWljcy5DUk0uY29ubmVjdGlvbnJlZmVyZW5jZSIsImNvbm5lY3Rpb25SZWZlcmVuY2VMb2dpY2FsTmFtZSI6ImF1dG9jb2Vfc2hhcmVkb2ZmaWNlMzY1dXNlcnNfMGI1YzAiLCJjb25uZWN0aW9uUmVmZXJlbmNlRGlzcGxheU5hbWUiOiJPZmZpY2UgMzY1IFVzZXJzIEF1dG9tYXRpb25Db0VDb3JlTWFpblByaXZhdGVQcmV2aWV3LTBiNWMwIiwiZGVzY3JpcHRpb24iOm51bGwsImNvbm5lY3RvcklkIjoiL3Byb3ZpZGVycy9NaWNyb3NvZnQuUG93ZXJBcHBzL2FwaXMvc2hhcmVkX29mZmljZTM2NXVzZXJzIiwiY29ubmVjdGlvbklkIjoiMWMzZjQ3OTEwMmZkNDcwZThiNzllNjU5OGIzYWVmMTIifSx7IkBvZGF0YS50eXBlIjoiTWljcm9zb2Z0LkR5bmFtaWNzLkNSTS5jb25uZWN0aW9ucmVmZXJlbmNlIiwiY29ubmVjdGlvblJlZmVyZW5jZUxvZ2ljYWxOYW1lIjoiYXV0b2NvZV9zaGFyZWRvZmZpY2UzNjVfMDEzMTMiLCJjb25uZWN0aW9uUmVmZXJlbmNlRGlzcGxheU5hbWUiOiJPZmZpY2UgMzY1IE91dGxvb2sgQXV0b21hdGlvbkNvRUNvcmVNYWluUHJpdmF0ZVByZXZpZXctMDEzMTMiLCJkZXNjcmlwdGlvbiI6bnVsbCwiY29ubmVjdG9ySWQiOiIvcHJvdmlkZXJzL01pY3Jvc29mdC5Qb3dlckFwcHMvYXBpcy9zaGFyZWRfb2ZmaWNlMzY1IiwiY29ubmVjdGlvbklkIjoiMWY3YzM3MzQ4OTNhNDk0NDk0NjZjODNhYmE0OWEwYTMifSx7IkBvZGF0YS50eXBlIjoiTWljcm9zb2Z0LkR5bmFtaWNzLkNSTS5jb25uZWN0aW9ucmVmZXJlbmNlIiwiY29ubmVjdGlvblJlZmVyZW5jZUxvZ2ljYWxOYW1lIjoiYXV0b2NvZV9zaGFyZWRwb3dlcnBsYXRmb3JtZm9yYWRtaW5zXzg5MGQ4IiwiY29ubmVjdGlvblJlZmVyZW5jZURpc3BsYXlOYW1lIjoiUG93ZXIgUGxhdGZvcm0gZm9yIEFkbWlucyBBdXRvbWF0aW9uQ29FTWFpbi04OTBkOCIsImRlc2NyaXB0aW9uIjpudWxsLCJjb25uZWN0b3JJZCI6Ii9wcm92aWRlcnMvTWljcm9zb2Z0LlBvd2VyQXBwcy9hcGlzL3NoYXJlZF9wb3dlcnBsYXRmb3JtZm9yYWRtaW5zIiwiY29ubmVjdGlvbklkIjoiNjk4MTU2ODg4NzI1NDUxYWJkMzFkODZjMjEwNGEzNjcifV0=|activateapprovalflow=True|activateroiflow=True|activatesyncflow=True|projectadminusers=|projectcontributors=|projectviewers=|businessowneremail=approver@contoso.com"
  package_settings = "installmainsolution=False|importconfigdata=True|activateapprovalflow=True|activateroiflow=True|activatesyncflow=True|projectadminusers=|projectcontributors=|projectviewers=|businessowneremail=approver@contoso.com"

  depends_on = [
    powerplatform_solution.creator_kit_solution_main,
    powerplatform_solution.automation_kit_main_solution
  ]
}

locals {
   empty_connection_reference_main_local = "${length(local_file.coe_automation_main_solution_settings.content) - 
   length(replace(local_file.coe_automation_main_solution_settings.content, "replace_with_your_connection_reference", ""))}"
}

resource "local_file" "coe_automation_main_solution_settings" {
  filename = "${path.module}/data/coe_automation_main_solution_settings.json"
  content = <<EOF
{
  "EnvironmentVariables": [
    {
      "SchemaName": "autocoe_DefaultFrequencyValues",
      "Value": "${var.env_variable_autocoe_default_frequency_values}"
    }
  ],
  "ConnectionReferences": [
    {
      "LogicalName": "autocoe_sharedapprovals_8865d",
      "ConnectionId": "${var.main_conn_ref_shared_approvals}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_approvals"
    },
    {
      "LogicalName": "autocoe_sharedapprovals_efae8",
      "ConnectionId": "${var.main_conn_ref_shared_approvals}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_approvals"
    },
    {
      "LogicalName": "autocoe_sharedcommondataserviceforapps_58896",
      "ConnectionId": "${var.main_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_sharedcommondataserviceforapps_9c637",
      "ConnectionId": "${var.main_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_sharedcommondataserviceforapps_b3698",
      "ConnectionId": "${var.main_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_sharedcommondataserviceforapps_b71a4",
      "ConnectionId": "${var.main_conn_ref_shared_commondataserviceforapps}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
    },
    {
      "LogicalName": "autocoe_sharedoffice365users_0b5c0",
      "ConnectionId": "${var.main_conn_ref_shared_office365users}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_office365users"
    },
    {
      "LogicalName": "autocoe_sharedoffice365users_15093",
      "ConnectionId": "${var.main_conn_ref_shared_office365users}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_office365users"
    },
    {
      "LogicalName": "autocoe_sharedoffice365_01313",
      "ConnectionId": "${var.main_conn_ref_shared_office365}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_office365"
    },
    {
      "LogicalName": "autocoe_sharedoffice365_5403c",
      "ConnectionId": "${var.main_conn_ref_shared_office365}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_office365"
    },
    {
      "LogicalName": "autocoe_sharedpowerplatformforadmins_5f3bd",
      "ConnectionId": "${var.main_conn_ref_shared_powerplatformforadmins}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_powerplatformforadmins"
    },
    {
      "LogicalName": "autocoe_sharedpowerplatformforadmins_890d8",
      "ConnectionId": "${var.main_conn_ref_shared_powerplatformforadmins}",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_powerplatformforadmins"
    }
  ]
}
EOF
}

