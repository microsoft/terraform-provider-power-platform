# Title

Insufficient Error Logging in `SendOperation`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

The error logging in the `SendOperation` function is insufficient when an error occurs during the API request execution. The error is passed to a conditional block for handling but is not consistently logged. If an error occurs and the `res` is `nil`, only the error message is stored in the `output` map, leaving crucial debugging and diagnostic information unavailable in the logs.

## Impact

- **Reduced debugging effectiveness**: Lack of detailed error logging makes debugging challenging, especially in production environments.
- **Poor observability**: Missing error context limits the ability to trace issues back to their source.
- **Inconsistent logging practices**: The function logs response details when successful but lacks adequate diagnostics when it encounters errors.

Severity: **Medium**

## Location

Found in `SendOperation`.

## Code Issue

```go
	if res == nil && err != nil {
		output["body"] = types.StringValue(err.Error())
	} else {
		if len(res.BodyAsBytes) > 0 {
			output["body"] = types.StringValue(string(res.BodyAsBytes))
		}
	}
```

## Fix

Introduce comprehensive error logging before storing the error message in the output map to improve error tracking and debugging.

```go
	if res == nil && err != nil {
		tflog.Error(ctx, fmt.Sprintf("API request failed: %v", err))
		output["body"] = types.StringValue(err.Error())
	} else {
		if len(res.BodyAsBytes) > 0 {
			output["body"] = types.StringValue(string(res.BodyAsBytes))
			tflog.Trace(ctx, fmt.Sprintf("API request succeeded: %v", string(res.BodyAsBytes)))
		}
	}
```

This improvement ensures errors are captured in the logs, making diagnostic efforts easier while maintaining consistent logging practices across the function.