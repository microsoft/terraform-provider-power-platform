# Title

Redundant Test Case Configurations

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

Many test cases repeat identical or minimally different configurations, leading to unnecessary redundancy. This makes the test file cumbersome to maintain and reduces readability.

## Impact

The redundancy increases maintenance costs and may lead to inadvertent errors during test updates. The severity of this issue is low, as it does not break functionality but impacts code efficiency and cleanliness.

## Location

Duplicate configurations are found in tests like `TestAccEnvironmentsResource_Validate_CreateGenerativeAiFeatures_US_Region_Update` and `TestAccEnvironmentsResource_Create_Environment_With_Env_Group`.

Example:
```go
resource "powerplatform_environment" "development" {
    display_name = "Example"
    location = "europe"
    environment_type = "Sandbox"
    // other identical attributes
}
```

## Fix

Extract common configurations into shared variables or utility functions and use these within the individual tests. This will reduce duplication and make the tests easier to manage.

```go
const commonEnvironmentConfig = `
resource "powerplatform_environment" "development" {
    display_name = "Example"
    location = "europe"
    environment_type = "Sandbox"
}`

func getEnvironmentConfig(additionalConfig string) string {
    return commonEnvironmentConfig + additionalConfig
}

// Example usage
{
    Config: getEnvironmentConfig(`
        dataverse {
            language_code = "1033"
            currency_code = "PLN"
        }
    `),
}
```

This approach maintains clarity, avoids repetition, and makes updates simpler.