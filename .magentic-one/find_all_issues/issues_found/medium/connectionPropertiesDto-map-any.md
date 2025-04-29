# Title

Incorrect use of `map[string]any` for JSON data.

## Path

/workspaces/terraform-provider-power-platform/internal/services/connection/dto.go

## Problem

In the `connectionPropertiesDto` struct, two fields (`ConnectionParameters` and `ConnectionParametersSet`) are declared as `map[string]any`. While using `interface{}` (or `any`, starting Go 1.18) provides flexibility, it can lead to runtime type issues when the JSON data structure is marshalled/unmarshalled, as no explicit type safety is guaranteed.

## Impact

This can lead to runtime type assertion failures or coding errors where unknown/unexpected types are encountered. Medium severity as it can cause runtime issues depending on data flows in production.

## Location

```go
connectionPropertiesDto struct {
    ConnectionParametersSet map[string]any  `json:"connectionParametersSet,omitempty"`
    ConnectionParameters    map[string]any  `json:"connectionParameters,omitempty"`
}
```

## Code Issue

```go
ConnectionParametersSet map[string]any  `json:"connectionParametersSet,omitempty"`
ConnectionParameters    map[string]any  `json:"connectionParameters,omitempty"`
```

## Fix

Define a struct for Connection Parameters, providing explicit types for the expected data JSON fields. Example:

```go
// Example struct for parameters, replace key/value types with expected ones:
type ConnectionParameter struct {
    Key   string      `json:"key"`
    Value interface{} `json:"value"`
}

connectionPropertiesDto struct {
    ConnectionParametersSet []ConnectionParameter `json:"connectionParametersSet,omitempty"`
    ConnectionParameters    []ConnectionParameter `json:"connectionParameters,omitempty"`
}

```

This ensures specific typing, helps avoid runtime bugs, and eliminates the use of `map[string]any`.