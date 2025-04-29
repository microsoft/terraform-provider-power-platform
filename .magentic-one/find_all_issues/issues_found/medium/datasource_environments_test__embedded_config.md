# Issue 5: Test Configuration Strings Directly Embedded in Code

## Path

`/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments_test.go`

## Problem

Test configurations, such as the `Config` field in resource `TestStep` definitions, are written directly in the code. For example:

```go
Config: `
resource "powerplatform_environment" "env" {
    display_name     = "` + mocks.TestName() + `"
    description      = "description"
    location         = "europe"
    azure_region     = "northeurope"
    environment_type = "Sandbox"
    cadence = "Moderate"
    dataverse = {
        language_code     = "1033"
        currency_code     = "USD"
        security_group_id = "00000000-0000-0000-0000-000000000000"
    }
}

data "powerplatform_environments" "all" {
    depends_on = [powerplatform_environment.env]
}

output "test_environment"{
    value = one([for env in data.powerplatform_environments.all.environments : env if env.id == powerplatform_environment.env.id])
}
`,
```

Embedding configuration in code reduces reusability, makes it harder to maintain, and complicates test debugging.

## Impact

Hardcoded test configurations reduce maintainability, especially for large projects, and make reusing test cases across different modules or testing environments challenging. Severity: **Medium**.

## Location

Found in the `TestAccEnvironmentsDataSource_Basic` and other test functions which embed configurations directly in their steps.

### Code Issue Example

```go
Steps: []resource.TestStep{
    {
        Config: `
        resource "powerplatform_environment" "env" {
            display_name     = "` + mocks.TestName() + `"
            description      = "description"
            location         = "europe"
            azure_region     = "northeurope"
            environment_type = "Sandbox"
            cadence = "Moderate"
            dataverse = {
                language_code     = "1033"
                currency_code     = "USD"
                security_group_id = "00000000-0000-0000-0000-000000000000"
            }
        }

        data "powerplatform_environments" "all" {
            depends_on = [powerplatform_environment.env]
        }

        output "test_environment"{
            value = one([for env in data.powerplatform_environments.all.environments : env if env.id == powerplatform_environment.env.id])
        }
        `,
    },
},
```

### Fix

Move the configuration to external files and load them dynamically:

1. Create a configuration file, e.g., `test_config.tf`.

```hcl
resource "powerplatform_environment" "env" {
    display_name     = "${var.test_name}"
    description      = "description"
    location         = "europe"
    azure_region     = "northeurope"
    environment_type = "Sandbox"
    cadence = "Moderate"
    dataverse = {
        language_code     = "1033"
        currency_code     = "USD"
        security_group_id = "00000000-0000-0000-0000-000000000000"
    }
}

data "powerplatform_environments" "all" {
    depends_on = [powerplatform_environment.env]
}

output "test_environment"{
    value = one([for env in data.powerplatform_environments.all.environments : env if env.id == powerplatform_environment.env.id])
}
```

2. Load this configuration file dynamically in the tests:

```go
config := helpers.LoadTestConfig("path/to/test_config.tf")
Steps: []resource.TestStep{
    {
        Config: config,
    },
}
```

This improves modularity and makes maintaining test configurations easier across the project.