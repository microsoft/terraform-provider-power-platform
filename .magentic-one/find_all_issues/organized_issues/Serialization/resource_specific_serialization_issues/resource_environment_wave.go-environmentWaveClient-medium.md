# Non-Idiomatic Client Struct Field Name

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go

## Problem

The field `EnvironmentWaveClient` in the `Resource` struct uses mixed casing (PascalCase) for a field that could be unexported. In Go, struct field names that do not need to be exported should start with a lowercase letter to indicate package-private scope and follow idiomatic Go conventions (unless the field is intended for serialization or use outside the current package).

## Impact

Non-idiomatic naming makes the codebase less consistent with Go conventions and could lead to confusion or accidental exposure of private details of a struct. Overall severity: **medium**.

## Location

```go
// Line (approximately 18)
type Resource struct {
	helpers.TypeInfo
	EnvironmentWaveClient *environmentWaveClient
}
```

## Code Issue

```go	type Resource struct {
	helpers.TypeInfo
	EnvironmentWaveClient *environmentWaveClient
}
```

## Fix

Make the field unexported by starting its name with a lowercase letter, unless it needs to be exported for reasons such as serialization, package usage, or framework requirements:

```go
type Resource struct {
	helpers.TypeInfo
	environmentWaveClient *environmentWaveClient
}
```

**Explanation:**
- This adheres to Go's idiomatic naming convention for struct fields.
- Only exported fields (uppercase) are accessible from outside the package.