# Struct Field Name Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/connection/models.go

## Problem

There is inconsistency in naming the fields representing Entra (Azure AD) object IDs. In some structs, it is called `EntraId`, while in others it is called `EntraObjectId`. Consistent naming should be maintained, especially as Go does not support field aliasing and this may cause confusion and access bugs.

## Impact

Medium. Confusing field names may cause bugs during struct usage, especially when mapping between similar models. Reduces code readability and increases maintenance burden.

## Location

- Structs:
  - `SharesPrincipalDataSourceModel` uses `EntraId`.
  - `SharePrincipalResourceModel` uses `EntraObjectId`.

## Code Issue

```go
type SharesPrincipalDataSourceModel struct {
	EntraId     types.String `tfsdk:"entra_object_id"`
	DisplayName types.String `tfsdk:"display_name"`
}

type SharePrincipalResourceModel struct {
	EntraObjectId types.String `tfsdk:"entra_object_id"`
	DisplayName   types.String `tfsdk:"display_name"`
}
```

## Fix

Choose a single convention (`EntraObjectId` is preferable for clarity) and use it everywhere:

```go
type SharesPrincipalDataSourceModel struct {
	EntraObjectId types.String `tfsdk:"entra_object_id"`
	DisplayName   types.String `tfsdk:"display_name"`
}

type SharePrincipalResourceModel struct {
	EntraObjectId types.String `tfsdk:"entra_object_id"`
	DisplayName   types.String `tfsdk:"display_name"`
}
```

This improves code consistency and reduces potential confusion and bugs.
