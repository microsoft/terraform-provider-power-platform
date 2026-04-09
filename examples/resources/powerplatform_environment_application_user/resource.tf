resource "powerplatform_environment_application_user" "example" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  application_id = "00000000-0000-0000-0000-000000000002"
  security_roles = [
    "MetaForm Global Admin",
    "MetaForm User",
  ]
}
