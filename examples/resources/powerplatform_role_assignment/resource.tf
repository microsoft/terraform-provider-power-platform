resource "powerplatform_role_assignment" "example" {
  environment_id     = "00000000-0000-0000-0000-000000000001"
  principal_id       = "00000000-0000-0000-0000-000000000002"
  security_role_name = "Basic User"
}
