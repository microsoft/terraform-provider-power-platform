# Title

Missing Null Check in `Configure` Method

##

`/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go`

## Problem

In the `Configure` function, while validating `req.ProviderData`, a null check is performed, but additional handling is required for the case when `req.ProviderData` is not null but contains unexpected types, which can lead to runtime panics.

## Impact

Incorrect assumption about the compatibility of the provided data could lead to runtime errors, affecting provider stability. Severity: **high**

## Location

Located in the `Configure` function of the `DataRecordDataSource` struct.

## Code Issue

```go
if req.ProviderData == nil {
	// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
	return
}

client := req.ProviderData.(*api.ProviderClient).Api

if client == nil {
	resp.Diagnostics.AddError(
		"Unexpected Resource Configure Type",
		fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	)

	return
}
```

## Fix

Add type validation before proceeding to use the `ProviderData`. This will reduce the likelihood of a runtime panic.

```go
if req.ProviderData == nil {
	// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
	return
}

// Check for proper type
providerClient, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
	resp.Diagnostics.AddError(
		"Unexpected Resource Configure Type",
		fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	)
	return
}

client := providerClient.Api

if client == nil {
	resp.Diagnostics.AddError(
		"Unexpected Resource Configure Type",
		"API client is nil. Please report this issue to the provider developers.",
	)

	return
}
```