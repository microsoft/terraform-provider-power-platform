# Title

Improper Regular Expression Usage

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

Regular expressions such as `regexp.MustCompile(helpers.GuidRegex)` are used without proper validation or checks in the test framework. Overreliance on regex matching without additional validation impacts code clarity and potential error handling.

## Impact

Improper regex can lead to test fragility. If the regex needs updates later, this code might fail silently, creating debugging challenges. Severity: Medium.

## Location

Example found at line 86, in the test configuration.

## Code Issue

```go
testMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", regexp.MustCompile(helpers.GuidRegex))
```

## Fix

Validate regexes prior to their implementation to ensure consistency and proactive error reporting.

```go
validatedRegex := helpers.ValidateRegex(helpers.GuidRegex) 
resource.TestMatchResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", validatedRegex)
```
