# Title

Duplicate or Outdated Test Data Source/Resource List Could Cause Fragile Tests

##

internal/provider/provider_test.go

## Problem

The test functions `TestUnitPowerPlatformProviderHasChildDataSources_Basic` and `TestUnitPowerPlatformProviderHasChildResources_Basic` hardcode lists of expected data sources/resources; if the provider registers new children or changes these methods, the tests may start failing due to code drift, or new features will not be automatically covered. Also, these lists can easily get out of sync with the real implementation, causing unnecessary maintenance overhead.

## Impact

Medium. More prone to errors when adding/changing data sources/resources, leading to fragile tests and more maintenance effort.

## Location

```go
expectedDataSources := []datasource.DataSource{
    analytics_data_export.NewAnalyticsExportDataSource(),
    powerapps.NewEnvironmentPowerAppsDataSource(),
    ...
}
expectedResources := []resource.Resource{
    environment.NewEnvironmentResource(),
    environment_groups.NewEnvironmentGroupResource(),
    ...
}
```

## Fix

Consider dynamically generating the lists or adding explicit notes for maintainers to update them with each change. At minimum, add comments/todos to help ensure these lists stay in sync, or add checks that compare only subsets, or detect changes with helpful output.

```go
// TODO: Keep this list in sync with provider.DataSources() registration
expectedDataSources := []datasource.DataSource{
    // ...
}
```

Or, if possible:

```go
// Optionally compare names or types, not raw instances, to allow more robust comparison.
```
