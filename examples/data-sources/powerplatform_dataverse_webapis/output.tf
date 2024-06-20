output "query_result" {
  description = "Query result"
  value       = jsondecode(data.powerplatform_dataverse_web_apis.webapi_query.output.body)
}
