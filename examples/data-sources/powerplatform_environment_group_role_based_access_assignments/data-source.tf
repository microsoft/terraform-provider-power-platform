# List all role assignments scoped to an environment group
data "powerplatform_environment_group_role_based_access_assignments" "example" {
  environment_group_id = "00000000-0000-0000-0000-000000000000"
}
