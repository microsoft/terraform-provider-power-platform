# Title

Excessive Inline Hardcoded Terraform Configs in Test Steps

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go

## Problem

In resource-based acceptance tests, the configs are built up as long Go strings with heavy inline string concatenation (especially with `mocks.TestName()` and similar). This makes configs harder to read, verify, and update, especially as they grow in size.

## Impact

This decreases readability, increases the difficulty of changing configs or referencing them elsewhere (like documentation or fixtures), and often leads to formatting errors. Severity: medium.

## Location

Check all `Config: ...` fields in resource test steps.

## Code Issue

```go
Config: `
  data "azuread_domains" "aad_domains" {
    only_initial = true
  }

  locals {
    domain_name = data.azuread_domains.aad_domains.domains[0].domain_name
  }

  resource "random_password" "passwords" {
      ...
  }

  resource "azuread_user" "test_user" {
    user_principal_name = "` + mocks.TestName() + `@${local.domain_name}"
    ...
  }

  resource "powerplatform_user" "new_user" {
    ...
  }`,
```

## Fix

Move such configs to multi-line raw string constants at the top of the test file, or (for larger providers) to fixtures files. Then, substitute in only the values that need to change at runtime.

```go
const baseUserConfig = `
data "azuread_domains" "aad_domains" { ... }
locals { ... }
...
resource "azuread_user" "test_user" {
  user_principal_name = "%s@${local.domain_name}"
  ...
}
...
`

userConfig := fmt.Sprintf(baseUserConfig, mocks.TestName())
```
