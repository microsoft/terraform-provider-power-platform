# Title

`solutionSettingsDto` should be exported

##

/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go

## Problem

The struct `solutionSettingsDto` is defined with a lowercase initial letter, making it unexported. If this struct is intended to be used outside the `solution` package (which is common for DTOs—Data Transfer Objects—especially in APIs or libraries), it should have an exported name (start with uppercase).

## Impact

Prevents use of the struct outside the package, which may cause maintainability or API usability issues. Severity: medium.

## Location

Line 6

## Code Issue

```go
type solutionSettingsDto struct {
	EnvironmentVariables []settingsEnvironmentVariableDto  `json:"environmentvariables"`
	ConnectionReferences []settingsConnectionReferencesDto `json:"connectionreferences"`
}
```

## Fix

Change the struct name to use an uppercase first letter:

```go
type SolutionSettingsDto struct {
	EnvironmentVariables []settingsEnvironmentVariableDto  `json:"environmentvariables"`
	ConnectionReferences []settingsConnectionReferencesDto `json:"connectionreferences"`
}
```

Also consider making the field types exported if they are intended for use outside the package as well.
