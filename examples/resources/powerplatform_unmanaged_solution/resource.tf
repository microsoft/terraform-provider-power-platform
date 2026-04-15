resource "powerplatform_unmanaged_solution" "solution" {
  environment_id = var.environment_id
  uniquename     = "TerraformUnmanagedSolution"
  display_name   = "Terraform Unmanaged Solution"
  publisher_id   = var.publisher_id
  description    = "Unmanaged solution created directly through the Dataverse solutions table."
}
