# Ensure a service principal exists as an application user with System Administrator role 
# in an imported environment
resource "powerplatform_environment_application_admin" "import_fix" {
  environment_id = "00000000-0000-0000-0000-000000000000" # GUID of environment
  application_id = "00000000-0000-0000-0000-000000000000" # GUID (client ID) of the service principal
}