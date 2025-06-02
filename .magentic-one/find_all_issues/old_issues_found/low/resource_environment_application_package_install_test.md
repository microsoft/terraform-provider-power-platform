### Issue Analysis:

# Title

Improve Inline Resource Configurations for Readability

##

/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install_test.go

## Problem

Inline configurations for Terraform resources and data sources are difficult to read and maintain. For example:

```go
Config: `
    resource "powerplatform_environment" "environment" {
        display_name = "` + mocks.TestName() + `"
        location = "unitedstates"
        environment_type = "Sandbox"
        ...
    }`,
```

The code combines inline variables and hardcoded values within a string literal, making readability cumbersome.

## Impact

This practice can lead to:

- Increased difficulty in debugging the configuration, especially if the inline variables are restructured.
- Reduced code clarity, which may make the integration and management more tedious for new developers.
  
Severity is **low** due to its minor impact on the overall functionality but gradually affects efficiency and maintainability.

## Location

Example location:

```go
Config: `
resource ....`
{

    return impact/!
.....bscleansgs