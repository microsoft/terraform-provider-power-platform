terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
  }
}

provider "powerplatform" {
  username = "${var.username}"
  password = "${var.password}"
  host = "http://localhost:8080"
}


resource "powerplatform_data_loss_prevention_policy" "my_policy" {
    display_name = "Test Policy"
    
    environment_type = "ExceptEnvironments"
    environment {
          name = "11111111-2222-3333-4444-555555555555"
    }
    environment {
          name = "Default-11111111-2222-3333-4444-555555555555"
    }

    default_connectors_classification = "Blocked"

    connector_group {
      classification = "Confidential"
      connector {
        id = "/providers/Microsoft.PowerApps/apis/shared_sql"
        name = "SQL Server"
      }
      connector {
        id = "/providers/Microsoft.PowerApps/apis/shared_assistantstudio"
        name = "Dynamics 365 Sales Insights"
      }
      connector {
        id = "/providers/Microsoft.PowerApps/apis/shared_dynamics365marketing"
        name = "Dynamics 365 Marketing"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_office365users"
        name = "Office 365 Users"
      }

    }

    connector_group {
      classification = "General"
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_office365"
        name = "Office 365 Outlook"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"
        name = "SharePoint"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_onedriveforbusiness"
        name = "OneDrive for Business"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_approvals"
        name = "Approvals"

      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_cloudappsecurity"
        name = "Defender for Cloud Apps"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_commondataservice"
        name = "Microsoft Dataverse (legacy)"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_commondataserviceforapps"
        name = "Microsoft Dataverse"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_excelonlinebusiness"
        name = "Excel Online (Business)"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_flowpush"
        name = "Notifications"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_kaizala"
        name = "Microsoft Kaizala"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_microsoftformspro"
        name = "Dynamics 365 Customer Voice"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_office365groups"
        name = "Office 365 Groups"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_office365groupsmail"
        name = "Office 365 Groups Mail"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_onenote"
        name = "OneNote (Business)"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_planner"
        name = "Planner"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_powerappsnotification"
        name = "Power Apps Notification"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_powerappsnotificationv2"
        name = "Power Apps Notification V2"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_powerbi"
        name = "Power BI"

      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_shifts"
        name = "Shifts for Microsoft Teams"

      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_skypeforbiz"
        name = "Skype for Business Online"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_teams"
        name = "Microsoft Teams"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_todo"
        name = "Microsoft To-Do (Business)"
      }
      connector {
        id =  "/providers/Microsoft.PowerApps/apis/shared_yammer"
        name = "Yammer"
      }
    }
    connector_group {
        classification = "Blocked"
    }
}

output "name" {
  value = powerplatform_data_loss_prevention_policy.my_policy.name
}

output "display_name"{
  value = powerplatform_data_loss_prevention_policy.my_policy.display_name
}