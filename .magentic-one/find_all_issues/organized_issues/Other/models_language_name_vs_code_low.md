# Incorrect Tag Name: `language_name` vs. Field Name `LanguageCode`

##

/workspaces/terraform-provider-power-platform/internal/services/environment/models.go

## Problem

In the `DataverseSourceModel` struct, the field is named `LanguageName`, but the tag is `tfsdk:"language_code"` and the type is `types.Int64`. This creates confusion, since all references and DTOs seem to operate with "language code", not "language name".

## Impact

- **Severity:** Low
- Reduces code clarity and consistency.
- May be confusing when mapping data to and from external sources/APIs, or for future maintainers.

## Location

```go
type DataverseSourceModel struct {
	...
	LanguageName        types.Int64  `tfsdk:"language_code"`
	...
}
```

## Code Issue

```go
LanguageName        types.Int64  `tfsdk:"language_code"`
```

## Fix

Rename the struct field to match the tag and usage everywhere in the codebase:

```go
LanguageCode        types.Int64  `tfsdk:"language_code"`
```

Be sure to update all logic that references `LanguageName` to now reference `LanguageCode`.
