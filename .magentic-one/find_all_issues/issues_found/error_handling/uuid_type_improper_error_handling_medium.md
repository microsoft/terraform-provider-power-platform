# Title

Improper Error Handling When Converting Diagnostics to Error

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_type.go

## Problem

In the method `ValueFromTerraform`, diagnostics (`diag.Diagnostics`) returned from `ValueFromString` is only checked for the presence of an error (`diags.HasError()`), and if true, it is converted to an error using `fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)`. This results in a loss of granular diagnostic information (for example, warnings or multiple error entries). Additionally, only the error information is carried forward, whereas the diagnostic data might be valuable for downstream consumers.

## Impact

Severity: **Medium**

- Loss of detailed diagnostics—only a generic error message is returned.
- Downstream consumers receiving the error do not get structured diagnostics, which can reduce the ability to provide feedback to users.
- Reduces overall debuggability and observability of error handling logic.

## Location

`func (t UUIDType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error)`  
Lines:

```go
	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}
```

## Code Issue

```go
	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}
```

## Fix

- You can propagate diagnostics or include more context, potentially by returning diagnostics (not just `error`) if your framework supports it.
- If only an `error` can be returned, you may want to at least extract and aggregate all errors from diagnostics, rather than just the full object as a formatted string.

Example:

```go
	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		var errMsg string
		for _, diag := range diags {
			errMsg += diag.Summary + ": " + diag.Detail + "; "
		}
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %s", errMsg)
	}
```

This makes the error message more useful by including all summaries and details from the diagnostics slice.

