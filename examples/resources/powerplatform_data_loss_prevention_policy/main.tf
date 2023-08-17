terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "github.com/microsoft/terraform-provider-power-platform"
    }
  }
}
provider "powerplatform" {
  username = var.username
  password = var.password
  tenant_id = var.tenant_id
}
resource "powerplatform_data_loss_prevention_policy" "my_policy" {
  display_name                      = "Test Policy"
  environment_type                  = "ExceptEnvironments"
  default_connectors_classification = "Blocked"
  environments = [
    {
      environment_name = "00000000-0000-0000-0000-000000000000"
    },
    {
      environment_name = "00000000-0000-0000-0000-000000000001"
    }
  ]

  connector_groups = [
    {
      classification = "Confidential"
      connectors = [{
        id   = "/providers/Microsoft.PowerApps/apis/shared_sql"
        name = "SQL Server"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_office365"
          name = "Office 365 Outlook"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_assistantstudio"
          name = "Dynamics 365 Sales Insights"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_dynamics365marketing"
          name = "Dynamics 365 Marketing"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_office365users"
          name = "Office 365 Users"
      }]
    },
    {
      classification = "General"
      connectors = [

       
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_onedriveforbusiness"
          name = "OneDrive for Business"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_approvals"
          name = "Approvals"

        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_cloudappsecurity"
          name = "Defender for Cloud Apps"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
          name = "Microsoft Dataverse (legacy)"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
          name = "Microsoft Dataverse"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_excelonlinebusiness"
          name = "Excel Online (Business)"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_flowpush"
          name = "Notifications"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_kaizala"
          name = "Microsoft Kaizala"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_microsoftformspro"
          name = "Dynamics 365 Customer Voice"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_office365groups"
          name = "Office 365 Groups"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_office365groupsmail"
          name = "Office 365 Groups Mail"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_onenote"
          name = "OneNote (Business)"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_planner"
          name = "Planner"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_powerappsnotification"
          name = "Power Apps Notification"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_powerappsnotificationv2"
          name = "Power Apps Notification V2"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_powerbi"
          name = "Power BI"

        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_shifts"
          name = "Shifts for Microsoft Teams"

        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_skypeforbiz"
          name = "Skype for Business Online"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_teams"
          name = "Microsoft Teams"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_todo"
          name = "Microsoft To-Do (Business)"
        },
        {
          id   = "/providers/Microsoft.PowerApps/apis/shared_yammer"
          name = "Yammer"
        },
         {
          id   = "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"
          name = "SharePoint"
        }
      ]
    },
    {
      classification = "Blocked"
      connectors     = []
    }
  ]

}
