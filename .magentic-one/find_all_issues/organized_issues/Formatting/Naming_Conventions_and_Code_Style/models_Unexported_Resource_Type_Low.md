# Issue 2: Unexported Type 'Resource' Used as Top-Level Type

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem

The struct `Resource` is unexported but defined at the top level in a file likely meant to describe schema/model information. In Go, exported public API structs/types should have an uppercase name so they are visible to consumers outside the package. If `Resource` is meant to be a package-level or public-facing type, it should be named `CopilotStudioAppInsightsResource` or similar, and exported. If not, re-evaluate its visibility and file location.

## Impact

Severity: **Low**

Could result in confusion or issues with maintainability, especially if the struct is intended to be used outside this package. It also represents a potential API misdesign.

## Location

```go
type Resource struct {
	helpers.TypeInfo
	CopilotStudioApplicationInsightsClient client
}
```

## Fix

Export the type if it's intended for public use, or move/rename it if it is internal-only.

```go
type CopilotStudioAppInsightsResource struct {
	helpers.TypeInfo
	CopilotStudioApplicationInsightsClient client
}
```
