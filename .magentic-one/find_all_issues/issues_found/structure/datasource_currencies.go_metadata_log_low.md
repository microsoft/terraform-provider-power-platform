# Logging Could Expose Sensitive Information

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go

## Problem

In the `Metadata` method, the following log statement:

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

logs the fully qualified type name as debug output in each metadata call. While this is generally harmless, it could expose internal implementation details that may reveal naming conventions, platform-specific identifiers, or otherwise-unintended metadata in user debug logs.

## Impact

Sensitive or internal-only information could be unintentionally exposed to users or in log aggregation platforms, especially in regulated environments or when running with debug enabled in CI/CD. This represents a minor but valid security and maintainability concern.

**Severity:** Low

## Location

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

## Fix

Remove or sanitize unnecessary logging, or log only non-sensitive and user-expected information.

```go
// Optionally remove or rephrase for clarity:
tflog.Debug(ctx, "Set datasource metadata")
```
