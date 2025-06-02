# Fragile Use of fmt.Sprintf for Logging in tflog.Debug

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

The code repeatedly uses `tflog.Debug(ctx, fmt.Sprintf(...))` instead of using tflog's structured arguments. The `tflog` package in the HashiCorp SDK supports structured logging, so using plain strings via `fmt.Sprintf` loses the key/value structure, making logs harder to filter and process (especially in large, multi-concurrent environments).

## Impact

- Logs become less searchable and lack clarity for automated tooling.
- Might lose useful information in CI/CD or while debugging.
- Reduces maintainability of provider logs.
- Severity: **low**.

## Location

In multiple sections, example:

```go
tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", d.FullTypeName()))
// And
tflog.Debug(ctx, fmt.Sprintf("Query: %s", query))
tflog.Debug(ctx, fmt.Sprintf("Headers: %v", headers))
```

## Code Issue

```go
tflog.Debug(ctx, fmt.Sprintf("Query: %s", query))
```

## Fix

Use tflog's structured message formatting for better diagnostics, e.g.:

```go
tflog.Debug(ctx, "Query requested",
	map[string]interface{}{
		"odata_query": query,
	})
tflog.Debug(ctx, "Headers sent for OData request", map[string]interface{}{"headers": headers})
```

Update all such usages throughout the file for clarity and structure.
