# Hardcoded Magic Strings in Test Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports_test.go

## Problem

Several resource attribute names, resource paths, and test data are hardcoded as strings throughout the test, both in the test configuration and in attribute checks. Using magic strings hampers maintainability and increases the risk of typos or mismatches if the schema evolves.

## Impact

Medium. Maintenance cost rises, and silent test failures can occur if underlying resource or provider attribute names change. This risk is amplified in a plugin ecosystem where provider schemas may evolve.

## Location

```go
resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", regexp.MustCompile(helpers.GuidRegex))
```

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", regexp.MustCompile(helpers.GuidRegex)),
resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.source", regexp.MustCompile(helpers.StringRegex)),
// ...etc.
```

## Fix

Define constants for attribute names and resource paths in a dedicated section or file. This centralizes modifications and reduces risk when renaming or refactoring schema fields.

```go
const (
    dataSourceName = "data.powerplatform_analytics_data_exports.test"
    attrExportsID = "exports.0.id"
    attrExportsSource = "exports.0.source"
    // ...etc.
)

resource.TestMatchResourceAttr(dataSourceName, attrExportsID, regexp.MustCompile(helpers.GuidRegex)),
resource.TestMatchResourceAttr(dataSourceName, attrExportsSource, regexp.MustCompile(helpers.StringRegex)),
// ...etc.
```
