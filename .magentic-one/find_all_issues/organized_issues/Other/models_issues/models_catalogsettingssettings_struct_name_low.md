# Issue: Potentially Incorrect Struct Field Name

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/models.go

## Problem

The struct `CatalogSettingsSettings` has a repetitive name. Usually the convention would be `CatalogSettings`.

## Impact

Double naming reduces code readability and increases the likelihood of usage errors elsewhere in the code. Severity: **low**.

## Location

- Line 117

## Code Issue

```go
type CatalogSettingsSettings struct {
	PowerCatalogAudienceSetting types.String `tfsdk:"power_catalog_audience_setting"`
}
```

## Fix

Rename the struct to `CatalogSettings`:

```go
type CatalogSettings struct {
	PowerCatalogAudienceSetting types.String `tfsdk:"power_catalog_audience_setting"`
}
```
