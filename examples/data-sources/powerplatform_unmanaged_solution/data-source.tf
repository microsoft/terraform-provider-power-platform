data "powerplatform_unmanaged_solution" "example" {
  environment_id = powerplatform_environment.example.id
  uniquename     = "TerraformUnmanagedSolution"
}
