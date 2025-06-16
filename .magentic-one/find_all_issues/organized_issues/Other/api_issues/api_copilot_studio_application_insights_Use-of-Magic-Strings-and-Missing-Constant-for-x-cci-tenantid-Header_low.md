# Title

Use of Magic Strings and Missing Constant for `x-cci-tenantid` Header

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/api_copilot_studio_application_insights.go

## Problem

The custom HTTP header `"x-cci-tenantid"` is hard-coded in multiple places; if the header name changes in the future, changes must be made in several places.

## Impact

This reduces code maintainability and increases the risk of bugs if the header name changes, **low severity** but poor hygiene.

## Location

```go
http.Header{"x-cci-tenantid": {env.Properties.TenantId}}
```

## Fix

Define a package constant for the header key and use it everywhere.

```go
const XCCITenantIDHeader = "x-cci-tenantid"

// ...

http.Header{XCCITenantIDHeader: {env.Properties.TenantId}}
```
