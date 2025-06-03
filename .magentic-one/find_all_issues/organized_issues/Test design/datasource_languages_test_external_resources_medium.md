# Direct Use of External Resources in Tests

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages_test.go

## Problem

The acceptance test (`TestAccLanguagesDataSource_Validate_Read`) appears to depend on live API interaction, which may lead to flaky or non-repeatable tests.

## Impact

This reduces reliability, makes CI/CD slower/fragile (medium severity for builds).

## Location

```go
func TestAccLanguagesDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			...
		},
	})
}
```

## Code Issue

```go
ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
```

## Fix

Ensure all external calls are properly mocked in test scope or isolate tests requiring live resources. Clearly separate and mark tests as "integration" vs. "unit" and ensure CI only runs stable tests by default. Provide documentation comments indicating live system dependency.
