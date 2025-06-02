# Inconsistent Attribute Naming: Dataverse vs. Environment

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

Throughout the test cases, the term `dataverse` is used for a nested block and its attributes are frequently mixed with other resource attributes (e.g., `dataverse.language_code`, `dataverse.organization_id`). However, other attributes are named at the top level (like `location`, `display_name`).

While this may reflect the schema of the actual resource, within Go code and tests, inconsistent scoping and naming for subresource attributes can reduce clarity. In some config snippets, sometimes the attribute nesting/scope is ambiguous ("dataverse" vs "environment"), which could lead to confusion in implementation, documentation, or onboarding.

## Impact

- **Severity: Low**
- This is a minor maintainability and readability nuisance, and can lead to confusion for contributors/readers unfamiliar with schema conventions.
- Not technically incorrect, but expert guidance suggests naming should be predictable and consistently scoped.

## Location

Example (from many spots in file):

```go
resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.currency_code", "PLN"),
resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
```

## Code Issue

```go
resource "powerplatform_environment" "development" {
    display_name = "..."
    location = "..."
    dataverse = {
        language_code = "1033"
        currency_code = "PLN"
        // ...
    }
}
```

## Fix

Consider making a clarification in both test helper naming and documentation as to why "dataverse" is at this level. Optionally, apply a clear prefix everywhere or encapsulate subresource checks and config generation in helper functions to minimize ambiguity.

```go
// Helper/clarification
func checkDataverseAttr(attr, expected string) resource.TestCheckFunc {
    return resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse."+attr, expected)
}

// Usage
Check: resource.ComposeTestCheckFunc(
    resource.TestCheckResourceAttr("powerplatform_environment.development", "location", "europe"),
    checkDataverseAttr("currency_code", "PLN"),
)
```


Also consider documenting this distinction at the top of the test file, so that reviewers and new maintainers can tell top-level vs. nested resource properties at a glance.
