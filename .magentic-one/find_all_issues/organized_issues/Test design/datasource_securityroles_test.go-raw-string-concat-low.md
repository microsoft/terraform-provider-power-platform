# Title

Overuse of Raw String Concatenation for Terraform Configurations in Test Steps

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/datasource_securityroles_test.go

## Problem

The test cases use raw string concatenation and inline string literals to construct Terraform configuration blocks inside resource test steps, making them error-prone and hard to read or maintain. For instance, `display_name = "` + mocks.TestName() + `"` is concatenated directly into the configuration string. This impacts readability and increases the likelihood of unescaped or broken strings and test flakiness.

## Impact

This reduces maintainability and increases the risk of subtle errors from copy-paste or lack of proper escaping/validation. Severity is **low** as it is mostly a readability/maintainability issue but can result in hard-to-find bugs if not addressed.

## Location

Example from the beginning of the file, in the `TestAccSecurityDataSource_Validate_Read` test step’s `Config`:

## Code Issue

```go
Config: `
    resource "powerplatform_environment" "env" {
        display_name      = "` + mocks.TestName() + `"
        location          = "unitedstates"
        environment_type  = "Sandbox"
        dataverse = {
            language_code     = "1033"
            currency_code     = "USD"
            security_group_id = "00000000-0000-0000-0000-000000000000"
        }
    }

    data "powerplatform_security_roles" "all" {
        environment_id = powerplatform_environment.env.id
    }`,
```

## Fix

Move test configuration strings into constant variables at the top of the test file or use `fmt.Sprintf()` to format in variables. This will improve readability and prevent accidental string errors.

```go
const envConfigTmpl = `
resource "powerplatform_environment" "env" {
    display_name      = "%s"
    location          = "unitedstates"
    environment_type  = "Sandbox"
    dataverse = {
        language_code     = "1033"
        currency_code     = "USD"
        security_group_id = "00000000-0000-0000-0000-000000000000"
    }
}

data "powerplatform_security_roles" "all" {
    environment_id = powerplatform_environment.env.id
}
`

// ...in test step:
Config: fmt.Sprintf(envConfigTmpl, mocks.TestName()),
```

---

This issue will be saved as  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/datasource_securityroles_test.go-raw-string-concat-low.md`.

Continuing to check for additional issues.
