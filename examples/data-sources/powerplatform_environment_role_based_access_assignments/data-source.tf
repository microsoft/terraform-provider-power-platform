# List all role assignments scoped to an environment
data "powerplatform_environment_role_based_access_assignments" "example" {
  environment_id = "00000000-0000-0000-0000-000000000000"
}
