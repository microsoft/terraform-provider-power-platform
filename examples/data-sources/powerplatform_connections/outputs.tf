output "all_connections" {
  description = "All connections avaiable in Power Platform for a given environment"
  value       = data.powerplatform_connections.all_connections
}
