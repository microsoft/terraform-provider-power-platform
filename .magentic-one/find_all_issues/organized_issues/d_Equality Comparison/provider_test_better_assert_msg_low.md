# Title

Test Functions Lack Descriptive Failure Messages or Context

##

internal/provider/provider_test.go

## Problem

The `require.Equal` and `require.Contains` assertions check for equality or membership but the failure messages are generic, e.g., "There are an unexpected number of registered resources". Adding actual/expected values or context would make failures easier to debug.

## Impact

Low. Improved error context helps developers quickly understand test failures.

## Location

```go
require.Equal(t, len(expectedDataSources), len(datasources), "There are an unexpected number of registered data sources")
require.Contains(t, expectedDataSources, d(), "An unexpected data source was registered")
```

## Fix

Enhance failure messages with actual values:

```go
require.Equalf(t, len(expectedDataSources), len(datasources), "Expected %d data sources, got %d", len(expectedDataSources), len(datasources))
require.Containsf(t, expectedDataSources, d(), "Data source %+v was not expected", d())
```
