output "main_solution_output" {
  value = module.power_platform.coe_automation_main_solution_settings_empty_connection_reference_length > 0 ? "Some connection references for AutomatioCoeMain solution are empty. Please add them and run 'terraform apply', so that operation can be completed." : null
}

output "coe_automation_satelite_solution_settings_empty_connection_reference_length" {
  value = module.power_platform.coe_automation_satelite_solution_settings_empty_connection_reference_length > 0 ? "Some connection references for AutomatioCoeSatelite solution are empty. Please add them and run 'terraform apply', so that operation can be completed." : null
}