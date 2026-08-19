output "role_assignments" {
  description = "All tenant scoped role assignments"
  value       = data.powerplatform_role_assignments.all.role_assignments
}
