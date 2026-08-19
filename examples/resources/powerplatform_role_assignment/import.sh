# Tenant scope: the role assignment id on its own
terraform import powerplatform_role_assignment.example 00000000-0000-0000-0000-000000000000

# Environment scope
terraform import powerplatform_role_assignment.environment environments/00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000000

# Environment group scope
terraform import powerplatform_role_assignment.environment_group environmentGroups/00000000-0000-0000-0000-000000000002/00000000-0000-0000-0000-000000000000
