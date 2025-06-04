# Naming: Inconsistent function names for test helpers (`TestName` vs `TestsEntraLicesingGroupName`)

##

/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go

## Problem

The file defines test helper functions with inconsistent and unclear naming conventions, such as:
- `TestName()`
- `TestsEntraLicesingGroupName()`

Specifically, `TestsEntraLicesingGroupName()` contains a typo (`Licesing` instead of `Licensing`) and an awkward pluralization (`Tests` vs `Test`). Additionally, the function name does not clearly describe its purpose (e.g., returning a static group name).

## Impact

Inconsistent or misleading naming reduces code readability and discoverability. Typos and unclear semantic intent can make code harder to maintain, refactor, and understand for other developers. Severity: **low**

## Location

Definition of helper:
```go
func TestsEntraLicesingGroupName() string {
	return "pptestusers"
}
```

## Code Issue

```go
func TestsEntraLicesingGroupName() string {
	return "pptestusers"
}
```

## Fix

Rename the helper with the correct spelling and consistent naming practice. For example:

```go
func TestEntraLicensingGroupName() string {
	return "pptestusers"
}
```
or simply
```go
func EntraLicensingGroupName() string {
	return "pptestusers"
}
```
and update any call sites.
