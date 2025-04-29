# Title

Ambiguity in Timeout Configuration (`Timeouts`)

##

`/workspaces/terraform-provider-power-platform/internal/services/data_record/models.go`

## Problem

Multiple structs (`DataRecordListDataSourceModel`, `DataRecordResourceModel`) use a field named `Timeouts` with type `timeouts.Value`. These struct definitions lack clarity on how this timeout value is set and managed. Missing configuration documentation or validation may result in undesired behavior during execution.

## Impact

Without proper validation or documentation, timeout values may go unchecked, leading to performance bottlenecks or failed operations in the services. This can adversely affect service reliability and user experience. Severity: **Medium**

## Location

### Appears in following locations:

```go
// Location 1
Timeouts timeouts.Value `tfsdk:"timeouts"`

// Location 2
Timeouts timeouts.Value `tfsdk:"timeouts"`
```

Found in structs:
```go
type DataRecordListDataSourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type DataRecordResourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
```

## Code Issue

```go
Timeouts timeouts.Value `tfsdk:"timeouts"`
```

## Fix

Introduce a validation mechanism for the configuration and provide documentation for appropriate timeout values. Include fallbacks or default values to ensure robust handling.

### Adjusted Code Example:

```go
type DataRecordListDataSourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// Validation for timeout values
func ValidateTimeouts(timeout timeouts.Value) error {
	if timeout < MinimalTimeout {
		return fmt.Errorf("timeout value (%d) is below the minimal allowed value (%d)", timeout, MinimalTimeout)
	}
	if timeout > MaximumTimeout {
		return fmt.Errorf("timeout value (%d) exceeds maximum allowed value (%d)", timeout, MaximumTimeout)
	}
	return nil
}

type DataRecordResourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// Default timeout configuration
func (resource *DataRecordResourceModel) DefaultTimeoutSetting() timeouts.Value {
	if resource.Timeouts == 0 {
		resource.Timeouts = DefaultTimeout
	}
	return resource.Timeouts
}
```

This fix ensures timeout values are strictly validated and default values are applied when needed.

