# Test Function Names Could be Clearer

##

/workspaces/terraform-provider-power-platform/internal/api/client_test.go

## Problem

Some test function names do not follow a consistent or idiomatic Go style. For example, prefixing with `TestUnit` is redundant; Go encourages just `Test...` and perhaps structuring with underscores for different behavior (e.g., `TestApiClient_GetConfig`). The current style does not impact correctness but can slightly hinder discoverability and readability when the test suite grows.

## Impact

Low severity—primarily impacts aesthetics and test organization.

## Location

```go
func TestUnitApiClient_GetConfig(t *testing.T) {
// ...
}
```
And similarly named test functions.

## Code Issue

Unidiomatic function name style in test definitions.

## Fix

Rename test functions using standard Go naming conventions. For example:

```go
func TestApiClient_GetConfig(t *testing.T) {
// ...
}
```
