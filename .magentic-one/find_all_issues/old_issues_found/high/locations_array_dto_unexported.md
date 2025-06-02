# Title

Unexported struct `locationsArrayDto`

## Path

/workspaces/terraform-provider-power-platform/internal/services/locations/dto.go

## Problem

The struct `locationsArrayDto` is unexported, yet its fields are designed for JSON serialization, implying external usage. As such, the struct should be exported for usage outside of the package.

## Impact

Similar to the `locationDto` issue, this can limit the usability of the struct elsewhere in the codebase and result in serialization/deserialization issues. **Severity: High**

## Location

Line containing `type locationsArrayDto struct`.

## Code Issue

```go
type locationsArrayDto struct {
    ID         string             `json:"id"`
    Type       string             `json:"type"`
    Name       string             `json:"name"`
    Properties locationProperties `json:"properties"`
}
```

## Fix

Change the struct to follow uppercase naming conventions to make it exported.

```go
type LocationsArrayDto struct {
    ID         string             `json:"id"`
    Type       string             `json:"type"`
    Name       string             `json:"name"`
    Properties LocationProperties `json:"properties"`
}
```

Ensure all usages of this struct are updated accordingly.
