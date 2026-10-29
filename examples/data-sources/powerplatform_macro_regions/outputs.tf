output "all_macro_regions" {
  description = "All macro region geographies available to this tenant"
  value       = data.powerplatform_macro_regions.all_macro_regions.macro_regions
}

output "macro_region_ids" {
  description = "Identifiers that can be used as the macro_region of a powerplatform_environment"
  value       = [for macro_region in data.powerplatform_macro_regions.all_macro_regions.macro_regions : macro_region.macro_region_id]
}
