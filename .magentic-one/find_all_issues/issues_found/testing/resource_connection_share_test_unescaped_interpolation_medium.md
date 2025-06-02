# Title

Unescaped string interpolation in Terraform test fixture may cause bugs

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share_test.go

## Problem

The test configuration string for the acceptance test includes interpolation of Go variables directly into Terraform HCL strings. Some of these insertions are used in fields like `"user_principal_name"` and `"display_name"`, without escaping (for example: `user_principal_name = "` + mocks.TestName() + `@${local.domain_name}"`). If the values returned by the mock functions contain special characters or edge cases, the resulting HCL may become invalid or fail unexpectedly.

## Impact

- May cause test failures when the values are unexpected or contain quotes, backslashes, or other HCL special characters.
- Difficult to debug, as the error only occurs when a specific edge case is hit.
- Medium severity: could reduce reliability if mock output changes or is extended.

## Location

In the body of the acceptance test configuration string:

```go
user_principal_name = "` + mocks.TestName() + `@${local.domain_name}"
display_name        = "` + mocks.TestName() + `"
mail_nickname       = "` + mocks.TestName() + `"
display_name     = "` + mocks.TestName() + `"
...
display_name   = "OpenAI Connection ` + mocks.TestName() + `"
```

## Code Issue

```go
user_principal_name = "` + mocks.TestName() + `@${local.domain_name}"
display_name        = "` + mocks.TestName() + `"
...
display_name   = "OpenAI Connection ` + mocks.TestName() + `"
```

## Fix

Escape user-injected content properly. For example, you can use `strconv.Quote` to safely insert a string value:

```go
import "strconv"

userPrincipal := mocks.TestName()
...
Config: `
resource "azuread_user" "test_user" {
    user_principal_name = ` + strconv.Quote(userPrincipal + "@${local.domain_name}") + `
    display_name        = ` + strconv.Quote(userPrincipal) + `
    mail_nickname       = ` + strconv.Quote(userPrincipal) + `
    password            = random_password.passwords.result
    usage_location      = "US"
}
...
`,
```

This ensures the string is safely encoded and prevents accidental breaking of the test HCL due to special characters.

---

This markdown will be saved to `/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/resource_connection_share_test_unescaped_interpolation_medium.md`.
