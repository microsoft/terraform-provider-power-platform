# Title

Incorrect Field Naming - Api Should Be API

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/api_analytics_data_exports.go

## Problem

The struct field name `Api` in the `Client` struct and related references (e.g., `client.Api`) do not follow Go's convention for well-known abbreviations. `API` should be in all upper case to improve code clarity and maintain consistency.

## Impact

Deviation from idiomatic naming conventions reduces the legibility and professionalism of the code, especially as the codebase grows and is maintained by more developers. The impact is low, but it's best-practice to resolve. Severity: Low.

## Location

```go
type Client struct {
	Api       *api.Client
	TenantApi tenant.Client
}
```

## Code Issue

```go
	Api       *api.Client
```

## Fix

Rename `Api` to `API` and update all its references within the file.

```go
type Client struct {
	API       *api.Client
	TenantAPI tenant.Client
}
```

For every usage—such as `client.Api`—rename to `client.API`.

---

This file will be saved to:

```
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/api_analytics_data_exports_naming_low.md
```
