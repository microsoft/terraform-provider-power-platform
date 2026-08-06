output "role_assignments" {
  description = "All role assignments scoped to the environment group"
  value       = data.powerplatform_environment_group_role_based_access_assignments.example.role_assignments
}
