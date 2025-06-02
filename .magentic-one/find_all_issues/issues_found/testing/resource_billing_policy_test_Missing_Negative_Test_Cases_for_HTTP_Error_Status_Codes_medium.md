# Missing Negative Test Cases for HTTP Error Status Codes

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

The current tests only cover successful HTTP codes (200, 201, 204), but do not verify behavior for 4xx/5xx responses from the API. Negative testing is crucial to ensure proper error handling.

## Impact

- **Test Coverage**: Lower confidence in error path handling.
- **Quality Assurance**: Potential for undetected bugs in error handling logic.

**Severity: Medium**

## Location

All test functions using:

```go
httpmock.RegisterResponder(...)
```

## Code Issue

No configuration for error responses.

## Fix

Add additional `TestStep` sections with POST, GET, and DELETE returning 4xx/5xx codes, asserting that Terraform properly returns an error.

```go
httpmock.RegisterResponder("POST", policyURL(""), func(req *http.Request) (*http.Response, error) {
    return httpmock.NewStringResponse(http.StatusBadRequest, ""), nil
})

// In the step:
{
    Config: `...`,
    ExpectError: regexp.MustCompile("BadRequest"),
}
```
