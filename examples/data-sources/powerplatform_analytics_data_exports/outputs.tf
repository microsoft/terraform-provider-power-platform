output "export_id" {
  description = "The ID of the analytics data export"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].id
}

output "source" {
  description = "The source of the analytics data"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].source
}

output "environments" {
  description = "The environments configured for analytics data export"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].environments[*].environment_id
}

output "organization_ids" {
  description = "The organization IDs for the environments"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].environments[*].organization_id
}

output "status" {
  description = "The current status of all analytics data exports"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].status[*]
}

output "sink_details" {
  description = "The sink configuration details for all exports"
  value = {
    for idx, export in data.powerplatform_analytics_data_exports.example.exports : idx => {
      id              = export.sink.id
      type            = export.sink.type
      subscription_id = export.sink.subscription_id
      resource_group  = export.sink.resource_group_name
      resource_name   = export.sink.resource_name
    }
  }
}

output "package_names" {
  description = "The package names for all analytics data exports"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].package_name
}

output "scenarios" {
  description = "The list of scenarios covered by each analytics export"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].scenarios
}

output "resource_providers" {
  description = "The resource providers for all analytics data exports"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].resource_provider
}

output "ai_types" {
  description = "The AI types for all analytics data exports"
  value       = data.powerplatform_analytics_data_exports.example.exports[*].ai_type
}
