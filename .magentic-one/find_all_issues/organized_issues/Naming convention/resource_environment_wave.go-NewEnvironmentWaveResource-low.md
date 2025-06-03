# Function Naming Not Following Go Idiom

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go

## Problem

The function `NewEnvironmentWaveResource` follows a standard Go convention for constructor patterns, but the filename and internal project naming might make it redundant or verbose. While not necessarily incorrect and widely accepted, if the package is already named `environment_wave`, having `EnvironmentWave` in the function name can be redundant. Simpler patterns with `NewResource` or more explicit descriptions are used in broader idiomatic Go libraries when the context (package) is clear. If the constructor is exported, this is a very minor naming/nesting nit.

## Impact

The impact is minimal: this naming issue is of **low severity**. However, more concise function names increase readability when used frequently in codebases.

## Location

```go
func NewEnvironmentWaveResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_wave",
		},
	}
}
```

## Code Issue

```go
func NewEnvironmentWaveResource() resource.Resource { ... }
```

## Fix

If consistent with project conventions, consider renaming to:

```go
func NewResource() resource.Resource { ... }
```

Or, if disambiguation between many resource types is necessary, keep as is.

**Note:** This is a style/readability note and not a correctness or maintainability problem.