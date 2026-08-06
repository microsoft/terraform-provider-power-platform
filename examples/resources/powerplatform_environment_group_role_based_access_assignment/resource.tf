# Assign a role to a service principal at the environment group level
resource "powerplatform_environment_group_role_based_access_assignment" "example" {
  environment_group_id = "00000000-0000-0000-0000-000000000000"
  principal_object_id  = "00000000-0000-0000-0000-000000000000"
  principal_type       = "ApplicationUser"
  role_definition_id   = "00000000-0000-0000-0000-000000000000"
}
