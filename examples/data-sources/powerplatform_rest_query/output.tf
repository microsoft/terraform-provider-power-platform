output "query_result" {
  description = "Query result"
  value       = jsondecode(data.powerplatform_rest_query.webapi_query.output.body)
}
