# Title

Large Terraform test configs impact test readability and maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share_test.go

## Problem

The acceptance test inline configuration (`Config:`) is very large, containing multiple distinct blocks (data, resources, locals) all concatenated in a single raw string with embedded Go variable interpolation. This impairs code readability and maintainability; any change to one resource makes the whole string harder to reason about, and risks introduce accidental errors or duplication. It also makes refactoring or sharing configs between tests difficult.

## Impact

- Decreased maintainability and code review clarity as test configs grow or change.
- Difficult to apply consistent escaping/interpolation fixes.
- Increases the risk of accidental copy-paste bugs.
- Difficult to unit test or build partial/incremental tests.
- Severity: medium.

## Location

The multiline string under the first test step in `TestAccConnectionsShareResource_Validate_Create`:

```go
Config: `
    ... A very large multiline string with all resource blocks inlined and interpolated ...
`,
```

## Code Issue

```go
Config: `
    data "azuread_domains" "aad_domains" {
        ...
    }
    ...
    resource "azuread_user" "test_user" {
        ...
    }
    ...
    resource "powerplatform_environment" "env" {
        ...
    }
    ...
    resource "powerplatform_connection" "azure_openai_connection" {
        ...
    }
    ...
    resource "powerplatform_connection_share" "share_with_user1" {
        ...
    }
`,
```

## Fix

Split out the configuration into helper functions or use template files for different resources/sections. For example:

```go
func testAzureUserConfig(userName, domain string) string {
    return fmt.Sprintf(`
resource "azuread_user" "test_user" {
    user_principal_name = "%s@%s"
    display_name        = "%s"
    ...
}
`, userName, domain, userName)
}

// In your main test, compose blocks
Config: `
    ` + testAzureUserConfig(name, domain) + `
    ` + testEnvConfig(envName) + `
    ...
`
```

Or maintain test fixtures as external `.tf` files and load/parameterize them, which also helps with syntax highlighting and validation.

---

This issue will be saved to `/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_connection_share_test_large_config_medium.md`.
