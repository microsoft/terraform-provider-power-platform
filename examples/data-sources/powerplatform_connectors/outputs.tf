output "all_connectors" {
  description = "All connectors avaiable in Power Platform"
  value       = data.powerplatform_connectors.all_connectors
}
