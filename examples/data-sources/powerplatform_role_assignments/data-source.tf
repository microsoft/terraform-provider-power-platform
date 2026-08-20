# List all tenant scoped role assignments
data "powerplatform_role_assignments" "all" {
  scope_type = "tenant"
}
