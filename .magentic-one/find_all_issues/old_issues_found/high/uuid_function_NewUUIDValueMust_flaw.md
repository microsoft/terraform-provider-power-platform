# Title

Error Prone Function Naming and Behavior

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid.go

## Problem

The function `NewUUIDValueMust` uses ambiguous wording "Must" in the name, suggesting a critical operation. However, its behavior and error handling are unclear. It merely calls another function (`ValueUUID`) without proper checks on its result. This might lead to unexpected runtime behavior or errors that are difficult to trace.

## Impact

This lack of robust error handling may lead to confusing runtime diagnostics, especially if `ValueUUID` fails silently or returns unexpected values. The issue is **high**, as it might introduce silent bugs or obscure error paths in the execution.

## Location

Function declaration: `NewUUIDValueMust`

## Code Issue

```go
func NewUUIDValueMust(value string) (UUID, diag.Diagnostics) {
	return NewUUIDValue(value).ValueUUID()
}
```

## Fix

Improve error handling by explicitly validating the result of `ValueUUID` and documenting potential error cases.

```go
// NewUUIDValueMust creates a UUID object with the specified string value and ensures its correctness.
// Parameters:
//   - value: the string to initialize the UUID with.
// Returns:
//   - A UUID instance along with diagnostics defining any issues encountered.
// Behavior:
//   - If creation fails due to invalid inputs, appropriate diagnostics are returned.
func NewUUIDValueMust(value string) (UUID, diag.Diagnostics) {
	uuid, diagnostics := NewUUIDValue(value).ValueUUID()
	if diagnostics.HasError() {
		return uuid, diagnostics
	}
	return uuid, nil
}
```