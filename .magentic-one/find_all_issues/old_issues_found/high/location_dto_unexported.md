# Title

Unexported struct `locationDto`

## Path

/workspaces/terraform-provider-power-platform/internal/services/locations/dto.go

## Problem

The struct `locationDto` is unexported (its name starts with a lowercase letter), but it contains fields that are marked with JSON tags, suggesting that it should likely be exposed or used in an external context. If the struct is meant to be used outside of this package, it should be exported.

## Impact

This issue can cause difficulties when trying to serialize or utilize the struct outside of its package, making the code less clean and potentially causing runtime errors or confusion for developers. **Severity: High**

## Location

Line containing `type locationDto struct`.

## Code Issue

```go
type locationDto struct { 
    Value []locationsArrayDto `json:"value"` 
}
```

## Fix

Change the struct to be exported by using an uppercase naming convention for the struct.

```go
type LocationDto struct {
    Value []LocationsArrayDto `json:"value"`
}
```

This ensures it can be accessed externally when required. Make sure to also refactor any uses of this struct throughout the codebase to match the updated name.
