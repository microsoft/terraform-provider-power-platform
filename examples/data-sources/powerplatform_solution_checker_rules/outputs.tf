output "environment_details" {
  description = "Details of the environment being used"
  value = {
    name     = powerplatform_environment.example.display_name
    type     = powerplatform_environment.example.environment_type
    location = powerplatform_environment.example.location
    id       = powerplatform_environment.example.id
  }
}

output "rules" {
  description = "List of all solution checker rules"
  value       = data.powerplatform_solution_checker_rules.example.rules
}

output "rule_count" {
  description = "Total number of solution checker rules"
  value       = length(data.powerplatform_solution_checker_rules.example.rules)
}
