# Possible Omission of Context Use for Timeout Handling

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go

## Problem

Several model structs embed a `timeouts.Value` field but do not include any mechanism for context management or passing context further down to API/resource calls. The presence of the timeout value implies intent to manage timeouts, but not ensuring `context.Context` is plumbed through may lead to inefficient cancellation handling or ignored timeouts.

## Impact

Timeouts may not be respected throughout the resource/service logic, leading to hanging routines or unsuccessful cancellations. This can create negative user experience and resource leaks. Severity: **medium**.

## Location

```go
type SecurityRolesListDataSourceModel struct {
	Timeouts       timeouts.Value                `tfsdk:"timeouts"`
	...
}
type UserResourceModel struct {
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	...
}
```

## Code Issue

```go
Timeouts       timeouts.Value                `tfsdk:"timeouts"`
Timeouts          timeouts.Value `tfsdk:"timeouts"`
```

## Fix

Wherever these models are used, ensure that their `Timeouts` value is utilized in conjunction with a `context.Context` that propagates the timeout duration to downstream API or database calls. 

For example, add a utility method:

```go
func (m *UserResourceModel) ContextWithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, m.Timeouts.Read)
}
```

Then, use the resulting context in all relevant operations.
