# Title

Missing Nil Check for DataRecordClient in Read Method

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

The `d.DataRecordClient` is used in the `Read` method without checking if it is nil. Since this attribute is set in the `Configure` method, but could possibly not be set (e.g., if `Configure` didn't run, or failed, or external test/misuse), calling a method on a nil struct will cause the application to panic.

## Impact

If `DataRecordClient` is nil, the call to `SendOperation` will result in a runtime panic, which will crash the provider. This is a critical error in production and should have protection.

## Location

`Read` method of `DataverseWebApiDatasource`

## Code Issue

```go
	outputObjectType, err := d.DataRecordClient.SendOperation(ctx, &DataverseWebApiOperation{
		Scope:              state.Scope,
		Method:             state.Method,
		Url:                state.Url,
		Body:               state.Body,
		Headers:            state.Headers,
		ExpectedHttpStatus: state.ExpectedHttpStatus,
	})
```

## Fix

Add a nil check for `d.DataRecordClient` before it's used and fail gracefully with a proper diagnostic error if it is nil:

```go
	if d.DataRecordClient == nil {
		resp.Diagnostics.AddError(
			"Not Configured",
			"DataRecordClient is not configured. This may indicate that the provider is not correctly initialized.",
		)
		return
	}

	outputObjectType, err := d.DataRecordClient.SendOperation(ctx, &DataverseWebApiOperation{
		Scope:              state.Scope,
		Method:             state.Method,
		Url:                state.Url,
		Body:               state.Body,
		Headers:            state.Headers,
		ExpectedHttpStatus: state.ExpectedHttpStatus,
	})
```
This prevents a panic and provides a more informative error to the user.
