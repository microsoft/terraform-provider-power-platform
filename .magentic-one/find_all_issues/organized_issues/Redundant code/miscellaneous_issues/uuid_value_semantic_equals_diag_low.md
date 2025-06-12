# Diagnostics Not Returned on Error in StringSemanticEquals

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_value.go

## Problem

In the `StringSemanticEquals` method, if errors occur when parsing the old or new UUID values using `uuid.ParseUUID`, the diagnostics are populated by calling `diags.AddError`, but the method only checks if diagnostics exist with `diags.HasError()` and then returns `false, diags`. However, the flow does not stop parsing the new UUID if the old UUID check fails, potentially leading to misleading or redundant errors.

## Impact

**Low Severity**: Combining both errors together may confuse users, as the error with the old value is as significant as the new value and error reporting could be clearer if the function returned immediately after the first encountered error.

## Location

Method `StringSemanticEquals`

## Code Issue

```go
	oldUUID, err := uuid.ParseUUID(v.ValueString())
	if err != nil {
		diags.AddError("expected old value to be a valid UUID", err.Error())
	}

	newUUID, err := uuid.ParseUUID(newValue.ValueString())
	if err != nil {
		diags.AddError("expected new value to be a valid UUID", err.Error())
	}

	if diags.HasError() {
		return false, diags
	}

	return reflect.DeepEqual(oldUUID, newUUID), diags
```

## Fix

Return diagnostics immediately after encountering the first error (optional: for stricter validation/reporting), or document the intent to collect all parsing errors. Here’s an example of returning on the first error (recommended for clarity):

```go
	oldUUID, err := uuid.ParseUUID(v.ValueString())
	if err != nil {
		diags.AddError("expected old value to be a valid UUID", err.Error())
		return false, diags
	}

	newUUID, err := uuid.ParseUUID(newValue.ValueString())
	if err != nil {
		diags.AddError("expected new value to be a valid UUID", err.Error())
		return false, diags
	}

	return reflect.DeepEqual(oldUUID, newUUID), diags
```

Or, if keeping the current style, clarify with appropriate comments.
