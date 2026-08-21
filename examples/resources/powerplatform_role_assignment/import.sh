# Tenant scope: the import id is the role assignment id on its own
# 00000000-0000-0000-0000-000000000000 = role assignment id
terraform import powerplatform_role_assignment.example 00000000-0000-0000-0000-000000000000

# Environment scope: environments/{environment id}/{role assignment id}
# 00000000-0000-0000-0000-000000000001 = environment id
# 00000000-0000-0000-0000-000000000000 = role assignment id
terraform import powerplatform_role_assignment.environment environments/00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000000

# Environment group scope: environmentGroups/{environment group id}/{role assignment id}
# 00000000-0000-0000-0000-000000000002 = environment group id
# 00000000-0000-0000-0000-000000000000 = role assignment id
terraform import powerplatform_role_assignment.environment_group environmentGroups/00000000-0000-0000-0000-000000000002/00000000-0000-0000-0000-000000000000
