# Title

Inconsistent JSON Tagging: Required vs Optional Tags

##

`/workspaces/terraform-provider-power-platform/internal/services/environment/dto.go`

## Problem

Some JSON struct tags such as `description` in `EnvironmentPropertiesDto` are not specified as `omitempty`. This leads to these fields always being included in serialized JSON objects, even when nil or empty.

## Impact

Including unnecessary fields in serialized output results in increased payload size and can confuse external systems consuming these data structures. Proper tagging is crucial for optimal API behavior and performance.

**Severity:** Medium

## Location

Struct `EnvironmentPropertiesDto` contains the field `Description` with the tag:

```go
Description string `json:"description"`
```

## Code Issue

```go
type EnviromentPropertiesDto struct {
    // Other fields omitted for brevity
    Description               string                            `json:"description"`
}
```

## Fix

Add `omitempty` to the JSON tag for `Description`.

```go
type EnvironmentPropertiesDto struct {
    // Other fields omitted for brevity
    Description string `json:"description,omitempty"`
}
```