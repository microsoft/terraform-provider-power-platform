# Unexported Struct Naming

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/models.go

## Problem

The struct `sourceModel` is unexported (starts with lowercase) but is used for what seems to be a resource model (based on its fields and tags). In Go, if this struct is meant to be accessed outside this package (e.g., from a provider, resource or test), it must be exported (renamed to `SourceModel`).

## Impact

If the struct is meant to be used only within this file or package, there is no direct issue. However, if it is used or should be used outside (based on SDK or provider resource pattern), not exporting it is a common oversight and could limit extensibility or integration.
**Severity: Medium**

## Location

Line: near start of struct `type sourceModel`

## Code Issue

```go
type sourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
	SystemId      types.String   `tfsdk:"system_id"`
	PolicyType    types.String   `tfsdk:"policy_type"`
}
```

## Fix

If meant to be exported, capitalize the struct name. Otherwise, clarify usage.

```go
type SourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
	SystemId      types.String   `tfsdk:"system_id"`
	PolicyType    types.String   `tfsdk:"policy_type"`
}
```
