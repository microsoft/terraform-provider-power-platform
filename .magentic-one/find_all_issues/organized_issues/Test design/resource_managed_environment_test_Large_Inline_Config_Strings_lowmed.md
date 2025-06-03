# Title

Redundant or Repeated Large Inline Config Strings in Test Cases

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment_test.go

## Problem

Configuration strings for creation and update scenarios are repeated in many `resource.TestStep` structs, often with only small differences (changing one resource argument or value). This increases maintenance burden, introduces risk of inconsistencies, and makes the tests harder to read and update.

## Impact

Low-Medium. Readability and maintainability are impacted; changes in the resource definition or required arguments may require updating many places. Although tests still function correctly, duplication is error-prone.

## Location

Throughout the file, especially in large acceptance test sequences like:

```go
Config: `
resource "powerplatform_environment" "development" {
    display_name     = "` + mocks.TestName() + `"
    location         = "unitedstates"
    environment_type = "Sandbox"
    dataverse = {
        language_code    = "1033"
        currency_code    = "USD"
        security_group_id = "00000000-0000-0000-0000-000000000000"
    }
}
resource "powerplatform_managed_environment" "managed_development" {
    // ...
}
`,
// ...same repeated in each TestStep with minor differences...
```

## Code Issue

```go
Config: `
resource "powerplatform_environment" "development" {
    // ...
}
resource "powerplatform_managed_environment" "managed_development" {
    // ...
}
`,
// ...copied and pasted in each step, with changes...
```

## Fix

Refactor by extracting shared config string logic into functions that accept arguments for variations. For example:

```go
func managedEnvConfig(envVars ...string) string {
    return fmt.Sprintf(`
resource "powerplatform_environment" "development" {
    display_name     = "%s"
    location         = "%s"
    environment_type = "%s"
    dataverse = {
        language_code    = "%s"
        currency_code    = "%s"
        security_group_id = "%s"
    }
}
resource "powerplatform_managed_environment" "managed_development" {
    environment_id             = powerplatform_environment.development.id
    // ... and so on ...
}
`, envVars[0], envVars[1], envVars[2], envVars[3], envVars[4], envVars[5])
}

// Usage inside TestStep
Config: managedEnvConfig(mocks.TestName(), "unitedstates", ...)
```

This consolidates config logic and makes the test steps smaller and less error-prone.

