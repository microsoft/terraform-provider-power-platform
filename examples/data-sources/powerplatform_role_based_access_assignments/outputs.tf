output "role_assignments" {
  description = "All tenant scoped role assignments"
  value       = data.powerplatform_role_based_access_assignments.all.role_assignments
}
