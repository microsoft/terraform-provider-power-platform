# Title

Inconsistent Error Message in `ValidateAttribute`

##

Path: `/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_value.go`

## Problem

The `ValidateAttribute` method creates a diagnostic error message when the UUID is invalid. However, the constructed error message uses a hardcoded string `"Invalid UUID String Value"` and a formatted string for details. This approach does not adhere to a consistent error-building standard across the implementation, which could lead to difficulties in localization or modification.

## Impact

Inconsistent error message construction reduces the maintainability and scalability of the code, especially when standardizing diagnostic messages is required. The severity is **low** because it does not directly affect functionality but impacts long-term code quality.

## Location

File: `/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_value.go`

## Code Issue

```go
func (v UUIDValue) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	if _, err := uuid.ParseUUID(v.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			UUIDTypeErrorInvalidStringHeader,
			fmt.Sprintf(UUIDTypeErrorInvalidStringDetails, v.ValueString()),
		)

		return
	}
}
```

## Fix

To ensure consistency in error construction, use a helper method or follow a standardized error definition approach.

```go
func (v UUIDValue) ValidateAttribute(ctx context.Context, req xattr.ValidateAttributeRequest, resp *xattr.ValidateAttributeResponse) {
	if v.IsUnknown() || v.IsNull() {
		return
	}

	if _, err := uuid.ParseUUID(v.ValueString()); err != nil {
		// Standardized error construction using helper function
		resp.Diagnostics.AddAttributeError(
			req.Path,
			GenerateUUIDErrorHeader(),
			GenerateUUIDErrorDetails(v.ValueString()),
		)

		return
	}
}

// Example helper functions for standardizing error messages
func GenerateUUIDErrorHeader() string {
	return UUIDTypeErrorInvalidStringHeader
}

func GenerateUUIDErrorDetails(value string) string {
	return fmt.Sprintf(UUIDTypeErrorInvalidStringDetails, value)
}
```

### Summary of Changes:
1. Introduced helper functions to standardize error message construction.
2. Improved maintainability by creating reusable methods for error headers and details.
