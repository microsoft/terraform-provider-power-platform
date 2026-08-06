# Assign a role to a service principal at the tenant level
resource "powerplatform_role_based_access_assignment" "example" {
  principal_object_id = "00000000-0000-0000-0000-000000000000"
  principal_type      = "ApplicationUser"
  role_definition_id  = "00000000-0000-0000-0000-000000000000"
}
