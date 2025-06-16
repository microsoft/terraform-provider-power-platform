# Issue: Inconsistent Naming Convention for Struct Fields

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go

## Problem

In the `Client` struct, the field `Api` is named using a mixed-case abbreviation (`Api` instead of `API`). Also, the field `TenantApi` uses the same pattern. Go standards recommend that abbreviations in names should be uppercase ("API" instead of "Api") for improved clarity and convention alignment.

## Impact

Low. While the code will work, it does not follow idiomatic Go naming conventions, which can create confusion and reduce codebase consistency especially for new contributors or when integrating with other Go tools.

## Location

```go
type Client struct {
	Api       *api.Client
	TenantApi tenant.Client
}
```

## Code Issue

```go
type Client struct {
	Api       *api.Client
	TenantApi tenant.Client
}
```

## Fix

Rename the fields to use proper uppercase abbreviation for "API":

```go
type Client struct {
	API       *api.Client
	TenantAPI tenant.Client
}
```

Update all usages of these fields and their constructor accordingly.

---
