output "all_connectors" {
  description = "All connection avaiable in Power Platform for a given environment"
  value       = data.powerplatform_connections.all_connections
}
