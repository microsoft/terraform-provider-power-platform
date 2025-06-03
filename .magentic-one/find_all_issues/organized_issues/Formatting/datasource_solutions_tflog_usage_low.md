# Potential Inefficiency: Logging with Sprintf Instead of Key/Value in tflog

##

/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions.go

## Problem

In the `Metadata` function, `tflog.Debug` is called with a formatted string, rather than structured fields. The terraform-plugin-log library encourages structured logging using attributes/fields for better log query, filtering, and searching.

## Impact

**Low**. Purely affects logging, but may make logs less queryable and less useful for debugging in complex environments.

## Location

Line 38:

## Code Issue

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

## Fix

Use structured key/value logging:

```go
tflog.Debug(ctx, "METADATA", map[string]any{
	"type_name": resp.TypeName,
})
```

