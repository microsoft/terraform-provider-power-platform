# Type Safety: Untyped Use of Slice

##

/workspaces/terraform-provider-power-platform/internal/services/rest/models.go

## Problem

The field `Headers []DataverseWebApiOperationHeaderResource` uses a slice of a custom struct, but its definition (`DataverseWebApiOperationHeaderResource`) is not visible in this file. If it's defined elsewhere and used widely, lack of further validation or description here provides little type safety or documentation for maintainers.

## Impact

**Severity: Low**

Without context, this can easily lead to misunderstandings of what is expected for this field, especially if users are unaware of the struct definition.

## Location

```go
Headers            []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
```

## Code Issue

```go
Headers            []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
```

## Fix

If possible, document the expected shape of `DataverseWebApiOperationHeaderResource` within this file or, at minimum, via a doc comment. This will help with future maintainability.

```go
// Headers is a list of HTTP header key-value pairs for the Dataverse call.
Headers            []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
```

And ensure `DataverseWebApiOperationHeaderResource` is well-defined and documented.

