# Code Structure – Lack of Constructor for Struct Initialization

##

/workspaces/terraform-provider-power-platform/internal/helpers/typeinfo.go

## Problem

There is no constructor function (e.g., `NewTypeInfo`) provided to enforce invariants (such as non-empty `TypeName`). Relying on direct struct initialization can lead to incomplete or invalid object state and makes testing and validation harder. Using a constructor could also centralize validation logic.

## Impact

Medium. Missing a constructor makes it easier for developers to create invalid state. Introducing a constructor would promote consistency, easier input validation, and better code maintainability.

## Location

```go
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}
```

## Code Issue

Direct struct initialization:
```go
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}
```

## Fix

Introduce a constructor to enforce invariants and initialize all required fields properly.

```go
func NewTypeInfo(providerType, typeName string) (*TypeInfo, error) {
	if typeName == "" {
		return nil, fmt.Errorf("TypeName must not be empty")
	}
	return &TypeInfo{
		ProviderTypeName: providerType,
		TypeName:         typeName,
	}, nil
}
```

Developers would then use the constructor rather than initializing the struct directly.
