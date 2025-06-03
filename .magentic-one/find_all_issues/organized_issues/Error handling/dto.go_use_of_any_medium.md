# Title

Use of Unstructured `any` (interface{}) Reduces Type Safety

##

/workspaces/terraform-provider-power-platform/internal/services/connection/dto.go

## Problem

Fields such as `ConnectionParametersSet`, `ConnectionParameters`, `Principal`, and `Capabilities` are typed as `map[string]any` or `[]any`. This severely reduces type safety, as anything can be placed in those maps/arrays, making it harder to catch errors at compile time and making downstream usage error-prone.

## Impact

Medium. Use of `any` (Go 1.18+ alias for `interface{}`) makes code more error-prone and reduces code clarity, as the contents and expected structure are implicit and not compiler-checked.

## Location

Multiple locations, for example:

```go
ConnectionParametersSet map[string]any  `json:"connectionParametersSet,omitempty"`
ConnectionParameters    map[string]any  `json:"connectionParameters,omitempty"`
Capabilities            []any           `json:"capabilities"`
Principal               map[string]any  `json:"principal"`
```

## Code Issue

```go
ConnectionParametersSet map[string]any  `json:"connectionParametersSet,omitempty"`
ConnectionParameters    map[string]any  `json:"connectionParameters,omitempty"`
Capabilities            []any           `json:"capabilities"`
Principal               map[string]any  `json:"principal"`
```

## Fix

If possible, define new struct types to represent the expected structure for these fields. For example:

```go
type ConnectionParameter struct {
    /* define expected fields here */
}

ConnectionParametersSet map[string]ConnectionParameter  `json:"connectionParametersSet,omitempty"`
ConnectionParameters    map[string]ConnectionParameter  `json:"connectionParameters,omitempty"`

type Capability struct { /* ... */ }
Capabilities []Capability `json:"capabilities"`

// For Principal, replace with a concrete type if possible:
Principal PrincipalDto `json:"principal"`
```

If you must support arbitrary data, leave a struct field as `any` with a comment explaining why (e.g., for truly dynamic JSON blobs), but document expectations.
