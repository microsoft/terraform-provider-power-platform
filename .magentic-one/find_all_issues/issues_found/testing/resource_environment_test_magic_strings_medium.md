# Use of Magic Strings in Test Inputs and Checks

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

Many test cases use hardcoded "magic string" values (e.g., GUIDs, location names, attribute names, currency codes) both in resource config and expected test checks. These repeated literal values appear across a wide variety of configs and mocks: e.g. `"00000000-0000-0000-0000-000000000001"`, `"europe"`, `"Sandbox"`, `"PLN"`, and more.

This practice makes the test code more brittle, harder to update, and increases duplication. Typos or changes to these canonical values require edit sweeps across the full file.

## Impact

- **Severity: Medium**
- Small changes to test fixture require global updates.
- Unintentional inconsistencies can slip in, especially with stringly-typed values like GUIDs or region names.
- Duplicates config meanings and distracts from test intent.

## Location

Example snippets (dozens of locations):

```go
"00000000-0000-0000-0000-000000000001"
"europe"
"Sandbox"
"PLN"
resource.TestCheckResourceAttr("powerplatform_environment.development", "dataverse.unique_name", "00000000-0000-0000-0000-000000000001"),
```

## Code Issue

```go
resource "powerplatform_environment" "development" {
    display_name = "displayname"
    location = "europe"
    environment_type = "Sandbox"
    dataverse = {
        currency_code = "PLN"
        domain = "00000000-0000-0000-0000-000000000001"
    }
}
```

## Fix

Define all standard/canonical strings as constants at the top or in a shared test helper, such as:

```go
const (
    testGuid = "00000000-0000-0000-0000-000000000001"
    testLocation = "europe"
    testType = "Sandbox"
    testCurrency = "PLN"
)
```

Or in a package-level `var` or `const` block that encapsulates common test-data for schema, region, and account info. Refer to those constants in configs, checks, and expectation matchers. This "de-magic-strings" the test definitions, improving readability and maintainability.
