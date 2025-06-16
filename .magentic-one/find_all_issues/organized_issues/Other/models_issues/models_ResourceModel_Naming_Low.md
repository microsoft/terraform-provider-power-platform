# Issue 5: Potential Type Naming Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem

The type name `ResourceModel` is overly generic and not scoped to the domain. In a large codebase, generic type names can conflict or overlap, reducing clarity. Consider using a more descriptive name like `CopilotStudioAppInsightsResourceModel`.

## Impact

Severity: **Low**

Generic naming reduces maintainability and code clarity.

## Location

```go
type ResourceModel struct {
	// ...
}
```

## Fix

Use a domain-specific name.

```go
type CopilotStudioAppInsightsResourceModel struct {
	// ...
}
```
