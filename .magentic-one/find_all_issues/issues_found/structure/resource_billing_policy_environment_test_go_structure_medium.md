# Title

Embedded Expression in Terraform Resource Configuration String Reduces Readability and Maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_environment_test.go

## Problem

The test cases for acceptance testing are constructed via large interpolated Go strings, with test resource and data values injected via string concatenation (e.g., `...name = "power-platform-billing-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"...`). This makes the configuration brittle, hard to maintain, and error-prone. Using template files or multi-line string variables improves readability and test maintenance.

## Impact

- Medium: While the code functions, maintenance costs increase, and debugging or making changes to test resources is more complicated and error-prone.
- Makes it difficult to see, at a glance, what the Terraform config structure is in tests.
- Susceptible to subtle bugs due to missing delimiters, improper escaping, or logic changes.

## Location

AccTest and UpdateTest blocks, configuration strings inside:

```go
Config: `
    data "azapi_client_config" "current" {}
    resource "azapi_resource" "rg_example" {
        ...
        name = "power-platform-billing-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
    }
    ...
`
```

## Code Issue

```go
name      = "power-platform-billing-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
...
display_name     = "billing_policy_example_environment_1_` + mocks.TestName() + `"
```

## Fix

Replace concatenated Go string building with a format string with placeholders, or extract the configuration to a helper/template to improve readability. Example:

```go
config := fmt.Sprintf(`
data "azapi_client_config" "current" {}

resource "azapi_resource" "rg_example" {
    type      = "Microsoft.Resources/resourceGroups@2021-04-01"
    location  = "East US"
    name      = "power-platform-billing-%s%d"
}

...
`, testName, randSuffix)
```

Or move the configuration to a separate file or string constant/template, loading/testing it with dynamic values replaced via template processing.

---

# Title

Redundant HTTP Mock Registration Blocks in Unit Tests

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_environment_test.go

## Problem

In both `TestUnitBillingPolicyResourceEnvironment_Validate_Create` and `TestUnitBillingPolicyResourceEnvironment_Validate_Update`, the HTTP mock registration for GET/POST API endpoints is repeated (hardcoded URLs). This code can be abstracted and reused via a helper or setup function for clarity and maintainability.

## Impact

- Low: Slightly increases test file verbosity and makes test maintenance harder.
- Leads to code duplication and future inconsistency risk if endpoints change in some tests but not others.

## Location

Occurrences in both test functions, e.g.:

```go
httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments?api-version=2022-03-01-preview`, ...)
httpmock.RegisterResponder("POST", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments/remove?api-version=2022-03-01-preview`, ...)
httpmock.RegisterResponder("POST", `https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000000/environments/add?api-version=2022-03-01-preview`, ...)
```

## Code Issue

Code block duplicated in `TestUnitBillingPolicyResourceEnvironment_Validate_Create` and `TestUnitBillingPolicyResourceEnvironment_Validate_Update`.

## Fix

Refactor the mock registration logic into a helper:

```go
func registerCommonBillingPolicyEnvMocks() {
    httpmock.RegisterResponder("GET", ... , ...)
    httpmock.RegisterResponder("POST", ..., ...)
    httpmock.RegisterResponder("POST", ..., ...)
}
```

And call this helper in each test:

```go
registerCommonBillingPolicyEnvMocks()
```
