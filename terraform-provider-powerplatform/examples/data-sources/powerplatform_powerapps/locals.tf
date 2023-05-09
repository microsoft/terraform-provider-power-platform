locals {
  description = "This local is used to filter down the list of Power Apps to a specific Power App in a environment"
  only_specific_name = toset(
    [
      for each in data.powerplatform_powerapps.specific_environment.apps :
      each if each.display_name == "Example App Display Name"
  ])
}
