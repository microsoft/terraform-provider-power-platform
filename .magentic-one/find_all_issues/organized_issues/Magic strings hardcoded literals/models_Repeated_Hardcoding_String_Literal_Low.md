# Issue 4: Repeated Hardcoding of String Literal

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem

The `NetworkIsolation` field is hardcoded as `"PublicNetwork"` in the DTO creation. If `"PublicNetwork"` is a constant value used throughout the package, it would be better to declare it as a `const` for reusability and to avoid typos.

## Impact

Severity: **Low**

Hardcoded literals lead to technical debt and increased risk of typos, especially when used in multiple locations or subject to change.

## Location

```go
NetworkIsolation:            "PublicNetwork",
```

## Fix

Define a package-level constant for the value.

```go
const DefaultNetworkIsolation = "PublicNetwork"

// ...
NetworkIsolation: DefaultNetworkIsolation,
```
