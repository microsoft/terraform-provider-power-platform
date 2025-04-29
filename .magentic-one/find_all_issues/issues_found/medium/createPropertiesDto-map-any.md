# Title

Reuse of `map[string]any` leading to potential runtime type issues.

## Path

/workspaces/terraform-provider-power-platform/internal/services/connection/dto.go

## Problem

In the `createPropertiesDto` struct, the fields `ConnectionParameters` and `ConnectionParametersSet` are declared as `map[string]any`. While this provides flexibility, it introduces runtime type assertion risks during JSON marshalling/unmarshalling due to lack of explicit type safety.

## Impact

Similar to earlier analysis for `connectionPropertiesDto`, this can lead to runtime failures or unexpected behaviors during data processing. Severity is medium as such issues can cause disruptions in production environments if improper data types are encountered.

## Location

```go
createPropertiesDto struct {
    ConnectionParametersSet map[string]any       `json:"connectionParametersSet,omitempty"`
    ConnectionParameters    map[string]any       `json:"connectionParameters,omitempty"`
}
```

## Code Issue

```go
ConnectionParametersSet map[string]any       `json:"connectionParametersSet,omitempty"`
ConnectionParameters    map[string]any       `json:"connectionParameters,omitempty"`
```

## Fix

Similar to the suggested fix for `connectionPropertiesDto`, define an explicit type for connection parameters instead of using `map[string]any`. Example:

```go
// Explicit struct for connection parameters:
type ConnectionParameter struct {
    Key   string      `json:"key"`
    Value interface{} `json:"value"`
}

createPropertiesDto struct {
    ConnectionParametersSet []ConnectionParameter `json:"connectionParametersSet,omitempty"`
    ConnectionParameters    []ConnectionParameter `json:"connectionParameters,omitempty"`
}
```

This ensures type safety and avoids runtime type assertion risks with `map[string]any`.