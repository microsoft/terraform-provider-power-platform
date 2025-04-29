# Issue Title

Lack of Error Handling in Test Configuration

---

## File Location

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy_test.go

---

## Problem Description

The code lacks adequate error handling within the test configurations, specifically in `TestCase` setups using `ExpectError` or similar mechanisms. Without these provisions, tests may pass erroneously when an API returns unexpected results.

Current configuration:

```go
{
    Config: testAccDlpPolicyConfig_basic(),
    Check: resource.ComposeTestCheckFunc(
        resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.#", "20"),
    ),
},
```

## Severity Level

Medium

---

## Impact Analysis

- Test reliability is compromised.
- Increases vulnerability to silent failures.
- Makes debugging harder in case of test failures.

---

## Proposed Solution

Introduce appropriate error-handling checks in the test scenarios:

```go
{
    Config: testAccDlpPolicyConfig_basic(),
    ExpectError: func(err error) bool {
        return err != nil // Customize as per API behavior
    },
    Check: resource.ComposeTestCheckFunc(
        resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.#", "20"),
    ),
},
```