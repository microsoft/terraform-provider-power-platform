# Use of Hardcoded API URLs and Query Parameters

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

API URLs and parameters are hardcoded in test mocks, which leads to duplication and brittle tests if the endpoints change.

## Impact

- **Maintainability**: Updating endpoints requires changing in multiple places.
- **Readability**: Magic strings obscure meaning.

**Severity: Low**

## Location

Throughout the usage of `httpmock.RegisterResponder`:

```go
httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview", ...)
```

## Code Issue

```go
httpmock.RegisterResponder("GET", "https://api.powerplatform.com/licensing/billingPolicies/00000000-0000-0000-0000-000000000001?api-version=2022-03-01-preview", ...)
```

## Fix

Define constants for common URLs and API versions at the top of your test file. Use formatted strings or helper functions to construct them.

```go
const baseBillingPolicyURL = "https://api.powerplatform.com/licensing/BillingPolicies"
const apiVersion = "2022-03-01-preview"

// Helper
func policyURL(policyID string) string {
    return fmt.Sprintf("https://api.powerplatform.com/licensing/billingPolicies/%s?api-version=%s", policyID, apiVersion)
}

httpmock.RegisterResponder("GET", policyURL("00000000-0000-0000-0000-000000000001"), ...)
```
