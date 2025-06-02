# Title

Invalid Struct Field Tags in `TenantApplicationPackageDataSourceModel`

##

`/workspaces/terraform-provider-power-platform/internal/services/application/models.go`

## Problem

The field tags of the struct `TenantApplicationPackageDataSourceModel` such as `application_descprition`, `unique_name`, and others are inconsistent and may not follow a proper convention. The misspelled tag keys like `application_descprition` also create difficulty in serialization/deserialization processes.

## Impact

Incorrect field tags lead to misinterpretation by Terraform SDK and other related tools. This may result in runtime errors or improper data handling. Severity: **high**.

## Location

Field tags within `TenantApplicationPackageDataSourceModel`.

## Code Issue

```go
ApplicationDescprition types.String `tfsdk:"application_descprition"`
UniqueName             types.String `tfsdk:"unique_name"`
```

## Fix

Ensure the field tags follow proper conventions and correct any spelling errors:

```go
ApplicationDescription types.String `tfsdk:"application_description"`
UniqueName             types.String `tfsdk:"unique_name"`
```
