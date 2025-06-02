### Title

Unoptimized Set Operations

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

### Problem

Set operations for attributes such as `AllowedIpRangeForFirewall`, `AllowedServiceTagsForFirewall`, and `ReverseProxyIpAddresses` lack efficient validation or transformation mechanisms:

```go
value := strings.Join(helpers.SetToStringSlice(securitySourceModel.AllowedIpRangeForFirewall), ",")
```

Direct transformations do not verify duplicates or format integrity.

### Impact

Severity: Medium.
Such inefficiencies can cause redundant computation, degraded memory usage, and associated runtime penalties.

### Location

File: models.go
Function: convertFromEnvironmentSecuritySettings.
Block: `AllowedIpRangeForFirewall` and `ReverseProxyIpAddresses` computations.

### Code Issue

```go
value := strings.Join(helpers.SetToStringSlice(securitySourceModel.AllowedIpRangeForFirewall), ",")
```

### Fix

Normalize, deduplicate, and validate formatting.

```go
value := normalizeSetData(helpers.SetToStringSlice(securitySourceModel.AllowedIpRangeForFirewall))

func normalizeSetData(values []string) string {
    seen := make(map[string]bool)
    var normalized []string
    for _, v := range values {
        if !seen[v] {
            seen[v] = true
            normalized = append(normalized, v)
        }
    }
    return strings.Join(normalized, ",")
}
```

This prevents duplication and ensures accuracy during runtime calls.
