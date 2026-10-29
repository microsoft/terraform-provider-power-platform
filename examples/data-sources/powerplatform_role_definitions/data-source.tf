# List all available role definitions
data "powerplatform_role_definitions" "all" {
}

output "roles" {
  value = data.powerplatform_role_definitions.all.role_definitions
}
