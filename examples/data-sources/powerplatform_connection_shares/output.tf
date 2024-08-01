output "all_shared_connections" {
  description = "All shares for a given connection and a given environment"
  value       = data.powerplatform_connection_shares.all_shares
}
