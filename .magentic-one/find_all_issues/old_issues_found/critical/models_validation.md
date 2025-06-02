# Title

Missing Validation for `DataverseWebApiDatasourceModel` Struct Fields

##
/workspaces/terraform-provider-power-platform/internal/services/rest/models.go

## Problem

The `DataverseWebApiDatasourceModel` struct does not appear to include any validation logic for its fields (e.g., `Scope`, `Method`, `Url`, etc.). The absence of validation poses a risk as invalid input data could make its way into the system, potentially causing failures or unintended behavior.

## Impact

Without validation steps for these fields:
- Incorrect or malformed data could be inadvertently processed.
- System errors may arise due to unsupported field values, leading to unpredictable behavior.
- Code defensiveness is reduced, which negatively impacts maintainability and reliability.

Severity: **Critical**

## Location

File: `/internal/services/rest/models.go`  
Struct: `DataverseWebApiDatasourceModel`

## Code Issue

```go
type DataverseWebApiDatasourceModel struct {
	Timeouts           timeouts.Value                           `tfsdk:"timeouts"`
	Scope              types.String                             `tfsdk:"scope"`
	Method             types.String                             `tfsdk:"method"`
	Url                types.String                             `tfsdk:"url"`
	Body               types.String                             `tfsdk:"body"`
	ExpectedHttpStatus []int                                    `tfsdk:"expected_http_status"`
	Headers            []DataverseWebApiOperationHeaderResource `tfsdk:"headers"`
	Output             types.Object                             `tfsdk:"output"`
}
```

## Fix

Add validation logic for each field, ensuring that values conform to expected formats and constraints. For example:

```go
import (
	"errors"
	"regexp"
)

// Validate validates the DataverseWebApiDatasourceModel structure.
func (m *DataverseWebApiDatasourceModel) Validate() error {
	// Validate scope
	if m.Scope == "" {
		return errors.New("scope cannot be empty")
	}

	// Validate method
	if m.Method == "" {
		return errors.New("method cannot be empty")
	}

	validMethods := []string{"GET", "POST", "PUT", "DELETE"}
	if !contains(validMethods, m.Method.ValueString()) {
		return errors.New("unsupported HTTP method")
	}

	// Validate URL
	if !isValidURL(m.Url.ValueString()) {
		return errors.New("invalid URL")
	}

	// Further field validations
	// ...

	return nil
}

// Helper function to check valid URL.
func isValidURL(url string) bool {
	pattern := `^https?://[\w\-\.]+$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// Helper function to check item existence in a slice.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
```

Integrate the `Validate` method into the workflow, ensuring validation occurs at the time of instantiation or before processing the struct.
