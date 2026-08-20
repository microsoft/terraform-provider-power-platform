# The composite id names the table the principal lives in and the immutable role id, because a
# system user id and a team id are rows in different tables, and role names can be renamed or
# duplicated across business units.

# assignment to a user or application user
terraform import powerplatform_security_role_assignment.example "00000000-0000-0000-0000-000000000001/systemusers/00000000-0000-0000-0000-000000000002/00000000-0000-0000-0000-000000000004"

# assignment to a team
terraform import powerplatform_security_role_assignment.team_admin "00000000-0000-0000-0000-000000000001/teams/00000000-0000-0000-0000-000000000003/00000000-0000-0000-0000-000000000004"
