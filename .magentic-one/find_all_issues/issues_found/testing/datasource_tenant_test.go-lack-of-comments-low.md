# Absence of Comments Explaining Test Intent

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant_test.go

## Problem

The test functions and steps have no comments or documentation explaining the intent of the test, the specific configuration being tested, or rationale for the expected results. While the attribute checks are self-explanatory for seasoned developers, future contributors (or auditors unfamiliar with the business domain) may find it hard to understand why certain attribute values or patterns are expected and what system behavior is validated.

## Impact

**Severity: Low**

- Decreases readability and maintainability.
- Slower onboarding for new contributors and code reviewers.
- Harder to spot when tests need updating as the product evolves.

## Location

Applies to both test functions and their key blocks, e.g.:

```go
func TestUnitTenantDataSource_Validate_Read(t *testing.T) {
    ...
    Steps: []resource.TestStep{
        {
            Config: `
            data "powerplatform_tenant" "tenant" {}`,

            Check: resource.ComposeAggregateTestCheckFunc(
                // tests...
            ),
        },
    },
}
```

## Fix

Add doc-comments on top of each test function describing what is under test and why, and inline comments in check lists or as necessary:

```go
// TestUnitTenantDataSource_Validate_Read verifies that the tenant datasource returns all
// expected fields when the mock API returns a valid payload. It also checks the correct mapping of API fields to schema.
func TestUnitTenantDataSource_Validate_Read(t *testing.T) {
    ...
    Steps: []resource.TestStep{
        {
            Config: `
            data "powerplatform_tenant" "tenant" {}`,

            // Checks: ensure all returned values match the static mocked API payload
            Check: resource.ComposeAggregateTestCheckFunc(
                // ...
            ),
        },
    },
}
```
