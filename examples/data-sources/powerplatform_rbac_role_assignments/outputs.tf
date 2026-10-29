output "role_assignments" {
  description = "All tenant scoped role assignments"
  value       = data.powerplatform_rbac_role_assignments.all.role_assignments
}
