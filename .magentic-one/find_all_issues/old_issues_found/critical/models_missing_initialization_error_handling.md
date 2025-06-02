# Title

Missing error handling for `DataSource` struct initialization.

##

`/workspaces/terraform-provider-power-platform/internal/services/capacity/models.go`

## Problem

The `DataSource` struct contains a `CapacityClient` field, but there is no indication of how errors will be handled during its initialization or function calls. The absence of error handling can lead to unpredictable behavior and harder debugging.

## Impact

This can cause runtime errors if `CapacityClient` fails during operations. Severity: **critical**.

## Location

Line defining `DataSource`.

## Code Issue

### Problematic Code:

```go
type DataSource struct {
	helpers.TypeInfo
	CapacityClient client
}
```

## Fix

Introduce error handling or validation mechanisms during the initialization of `DataSource`. For example:

```go
type DataSource struct {
	helpers.TypeInfo
	CapacityClient client
}

// Initialize DataSource
func NewDataSource(capacityClient client) (*DataSource, error) {
	if capacityClient == nil {
		return nil, fmt.Errorf("CapacityClient cannot be nil")
	}
	return &DataSource{
		CapacityClient: capacityClient,
	}, nil
}
```