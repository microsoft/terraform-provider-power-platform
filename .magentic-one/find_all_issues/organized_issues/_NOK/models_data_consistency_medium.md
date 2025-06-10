# Data Consistency: Use of `types.String` for All Fields

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/models.go

## Problem

All fields in `EnvironmentPowerAppsDataSourceModel` are defined as `types.String`, including `CreatedTime`. If `CreatedTime` represents a timestamp or date, using a string type might lead to data parsing and validation issues.

## Impact

Improper data types can cause runtime errors, issues in parsing, and hinder validation/formatting of date/time fields (severity: medium).

## Location

```go
type EnvironmentPowerAppsDataSourceModel struct {
	EnvironmentId types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Name          types.String `tfsdk:"name"`
	CreatedTime   types.String `tfsdk:"created_time"`
}
```

## Code Issue

```go
	CreatedTime   types.String `tfsdk:"created_time"`
```

## Fix

Consider converting `CreatedTime` to a proper time or timestamp type (if supported by the framework, e.g., `types.Time`) and ensure that `ConvertFromPowerAppDto` properly sets this:

```go
// If framework supports types.Time:
import "github.com/hashicorp/terraform-plugin-framework/types"

type EnvironmentPowerAppsDataSourceModel struct {
	EnvironmentId types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Name          types.String `tfsdk:"name"`
	CreatedTime   types.Time   `tfsdk:"created_time"`
}
```

And update the conversion function to parse the time accordingly.
