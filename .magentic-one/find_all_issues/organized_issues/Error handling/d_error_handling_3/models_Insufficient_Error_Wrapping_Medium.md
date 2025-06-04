# Issue 3: Insufficient Error Wrapping in Data Conversion

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem

Both `createAppInsightsConfigDtoFromSourceModel` and `convertAppInsightsConfigModelFromDto` return a simple error without rich context. If errors are ever generated in these functions (such as after adding validation, or if ValueString/ValueBool might error in the future), they should be wrapped or annotated to provide meaningful call-site information.

## Impact

Severity: **Medium**

Lack of context in errors makes debugging and tracing issues more difficult in complex workflows.

## Location

```go
return nil, fmt.Errorf("EnvironmentId cannot be empty")
```

And similar.

## Fix

Use Go 1.13+ error wrapping for context where appropriate.

```go
return nil, fmt.Errorf("failed to create AppInsightsConfigDto: EnvironmentId cannot be empty")
```

And similarly for other errors.
