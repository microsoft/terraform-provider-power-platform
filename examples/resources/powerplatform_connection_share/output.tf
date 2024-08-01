output "sql_connection" {
  description = "New Sql connection"
  value       = powerplatform_connection.new_sql_connection
  sensitive   = true
}

output "connection_share" {
  description = "Share the connection with admin"
  value       = powerplatform_connection_share.share_with_admin
}
