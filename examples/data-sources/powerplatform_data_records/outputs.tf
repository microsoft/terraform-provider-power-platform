output "query_result" {
  description = "Query result"
  value       = data.powerplatform_data_records.data_query.rows
}
