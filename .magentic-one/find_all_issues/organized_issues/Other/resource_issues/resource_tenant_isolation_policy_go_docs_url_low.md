# Use of Deprecated or Outdated Documentation Links

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

The markdown documentation description of the schema includes the hard-coded link:  
`https://docs.microsoft.com/en-us/power-platform/admin/cross-tenant-restrictions`  
which is a legacy link. The preferred domain is now `learn.microsoft.com`, and keeping links current/maintained improves user experience and documentation accuracy.

## Impact

- **Severity: Low**
- Broken or outdated links can cause confusion or erode trust in the provider's documentation quality.
- Users may not have access to the most current documentation.

## Location

```go
MarkdownDescription: "Manages a Power Platform tenant isolation policy. ... Learn more: https://docs.microsoft.com/en-us/power-platform/admin/cross-tenant-restrictions",
```

## Code Issue

```go
MarkdownDescription: "Manages a Power Platform tenant isolation policy. Tenant isolation can be used to block external tenants " +
    "... Learn more: https://docs.microsoft.com/en-us/power-platform/admin/cross-tenant-restrictions",
```

## Fix

Update the URL to follow the latest Microsoft Docs conventions.  If such content exists at the new address, use:

```go
MarkdownDescription: "Manages a Power Platform tenant isolation policy. ... Learn more: https://learn.microsoft.com/en-us/power-platform/admin/cross-tenant-restrictions",
```
