# Title

Unmocked External API Calls

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

Several external API calls rely on hardcoded URLs and data. For example, in `httpmock.RegisterResponder`, the actual API endpoint is hardcoded when responding to requests. The solution should dynamically mock data reported by these requests.

## Impact

Hardcoded API-related calls can break tests in varied environments, lead to false positives, and heavily reduce test flexibility. Severity: High.

## Location

Lines 45, 48, 51.

## Code Issue

```go
httpmock.RegisterResponder("POST", "https://api.powerplatform.com/licensing/BillingPolicies?api-version=2022-03-01-preview"
```

## Fix

Replace hardcoded API endpoints with configurable mocks.

```go
httpmock.RegisterResponder("POST", mocks.MockAPIEndpoint("BillingPolicies")

where MockAPIEndpoint dynamically generates the API URL for test execution yet without false issues, maintaining ease+=