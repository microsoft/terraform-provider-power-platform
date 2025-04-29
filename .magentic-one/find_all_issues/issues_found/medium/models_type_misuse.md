# Title

Potential Type Misuse in `DataverseSourceModel`

## Path

`/workspaces/terraform-provider-power-platform/internal/services/environment/models.go`

## Problem

In the `DataverseSourceModel`, the `LanguageName` field is defined as `types.Int64`. However, this field might represent a language code, which is typically more suitable as a `string` for flexibility and readability, especially if different representations (e.g., "en-us", "1033") might be used.

## Impact

- Usage of `types.Int64` limits future enhancements that might involve non-numeric representations.
- Reduces code readability and understanding of the field's intent.
- Severity: Medium.

## Location

Defined in `DataverseSourceModel`.

## Code Issue

```go
type DataverseSourceModel struct {
	LanguageName        types.Int64  `tfsdk:"language_code"`
	// other fields...
}
```

## Fix

Change the type of `LanguageName` to `types.String` for better flexibility and semantic clarity.

```go
type DataverseSourceModel struct {
	// LanguageName refers to the code of the language in use, e.g., "en-us", "1033".
	LanguageName        types.String `tfsdk:"language_code"`
	// other fields...
}
```