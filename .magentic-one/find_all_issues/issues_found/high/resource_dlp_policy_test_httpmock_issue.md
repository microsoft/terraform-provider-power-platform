# Title

Repeated Re-Activation and Deactivation of `httpmock`

## Path

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

`httpmock.Activate()` and `httpmock.DeactivateAndReset()` are being executed repeatedly across multiple test cases, creating redundancy and increasing the maintenance overhead. This could also lead to potential errors if `DeactivateAndReset()` is not called properly in some cases.

## Impact

This issue increases the risk of inconsistent tests, potential resource leaks, and makes it harder to read or maintain the test suite. Severity is **high**, as it directly affects the reliability of tests.

## Location

Many instances found in `TestUnitDataLossPreventionPolicyResource_Validate_Update`, `TestUnitDataLossPreventionPolicyResource_Validate_Create`, and other test functions within this file.

## Code Issue

```go
func TestUnitDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies`,
        func(req *http.Request) (*http.Response, error) {
            return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create/get_policy_00000000-0000-0000-0000-000000000001.json").String()), nil
        })
    // Additional test logic
}
```

## Fix

Use a helper function or setup method that activates and deactivates `httpmock` universally for all tests, ensuring DRY (Don't Repeat Yourself) principles are applied.

### Example Fix

```go
func setupHttpMock() func() {
    httpmock.Activate()
    return func() {
        httpmock.DeactivateAndReset()
    }
}

func TestUnitDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
    tearDown := setupHttpMock()
    defer tearDown()

    httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies`,
        func(req *http.Request) (*http.Response, error) {
            return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create/get_policy_00000000-0000-0000-0000-000000000001.json").String()), nil
        })
    // Additional test logic
}
```