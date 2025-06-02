# Test Function Naming Not Consistent

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query_test.go

## Problem

The test functions in the file use both `TestUnitXxx` and `TestAccXxx` prefixes, which is good, but the `Using_Scope` part uses underscores, which is not the convention in Go test function naming. By convention, camel case is used. Using underscores diverts from Go idiomatic naming and can reduce readability.

## Impact

This is a **low** severity issue. It does not break any functionality but makes the test naming non-idiomatic, which may impact readability and consistency with other Go tests in the codebase.

## Location

```go
func TestUnitDatasourceRestQuery_WhoAmI_Using_Scope(t *testing.T) {
...
}

func TestAccDatasourceRestQuery_WhoAmI_Using_Scope(t *testing.T) {
...
}
```

## Code Issue

```go
func TestUnitDatasourceRestQuery_WhoAmI_Using_Scope(t *testing.T) {
	// ...
}

func TestAccDatasourceRestQuery_WhoAmI_Using_Scope(t *testing.T) {
	// ...
}
```

## Fix

Change function names to use camel case, and remove underscores for better Go style compliance.

```go
func TestUnitDatasourceRestQueryWhoAmIUsingScope(t *testing.T) {
	// ...
}

func TestAccDatasourceRestQueryWhoAmIUsingScope(t *testing.T) {
	// ...
}
```
