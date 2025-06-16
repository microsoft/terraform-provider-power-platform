# Naming: Struct Field Naming Convention for Acronyms

##

/workspaces/terraform-provider-power-platform/internal/services/rest/models.go

## Problem

Field names such as `Url` and `Http` in `DataverseWebApiDatasourceModel` do not use the Go convention for acronyms, which is `URL` and `HTTP` respectively. Go recommends all uppercase for acronyms to improve readability and consistency.

## Impact

**Severity: Low**

While not causing runtime issues, this diminishes consistency, may be confusing to readers, and does not align with Go idioms.

## Location

```go
type DataverseWebApiDatasourceModel struct {
    Timeouts           timeouts.Value                           `tfsdk:"timeouts"`
    Scope              types.String                             `tfsdk:"scope"`
    Method             types.String                             `tfsdk:"method"`
    Url                types.String                             `tfsdk:"url"`
    Body               types.String                             `tfsdk:"body"`
    ExpectedHttpStatus []int                                    `tfsdk:"expected_http_status"`
    Headers            []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
    Output             types.Object                             `tfsdk:"output"`
}
```

## Code Issue

```go
    Url                types.String                             `tfsdk:"url"`
    ExpectedHttpStatus []int                                    `tfsdk:"expected_http_status"`
```

## Fix

Follow Go conventions by using uppercase for acronyms in field names:

```go
    URL                types.String                             `tfsdk:"url"`
    ExpectedHTTPStatus []int                                    `tfsdk:"expected_http_status"`
```

Rename all usages accordingly throughout the codebase.

