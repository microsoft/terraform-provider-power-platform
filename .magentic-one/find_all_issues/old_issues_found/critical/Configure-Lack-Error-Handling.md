# Title

Lack of Error Handling for `newWebApiClient` initialization

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest.go

## Problem

In the `Configure` function, the `newWebApiClient(client.Api)` call does not perform error checking when initializing `DataRecordClient`. If the initialization fails, the program could proceed with an invalid client, resulting in runtime errors.

## Impact

This omission may lead to unpredictable behavior or failures when attempting operations that rely on the initialized client. Severity: **Critical**.

## Location

Function: `Configure`

## Code Issue

```go
	r.DataRecordClient = newWebApiClient(client.Api)
```

## Fix

Add error handling to validate the result of `newWebApiClient`. For example:

```go
	clientInstance, err := newWebApiClient(client.Api)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to initialize Web API Client",
			fmt.Sprintf("Error: %v", err),
		)
		return
	}
	r.DataRecordClient = clientInstance
```

This ensures that the client is valid before proceeding with resource configuration.