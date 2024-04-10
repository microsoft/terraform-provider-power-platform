output "all_settings" {
  description = "Environment settings for a chosen environment."
  value       = data.powerplatform_environment_settings.settings
}
