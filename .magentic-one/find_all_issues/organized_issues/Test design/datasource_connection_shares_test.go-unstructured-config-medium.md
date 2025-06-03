# Unstructured and Complex Inline Terraform Configuration in Tests

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares_test.go

## Problem

The acceptance test (`TestAccConnectionsShareDataSource_Validate_Read`) embeds a massive, deeply nested HCL (Terraform) configuration as a raw backtick string inline in the test code. This includes repeated interpolation and concatenation of results from mock functions. This approach makes it difficult to maintain, understand, and reason about the test cases. Any error in the inline HCL, improper escaping, or string formatting mistakes (like missing commas, typos, or missing quotes) can cause the test setup to silently misbehave or fail in confusing ways.

## Impact

Severity: **Medium**  
This hampers code readability and maintainability. It also increases the surface for copy-paste errors and makes reviewing, updating, or extending test scenarios more error-prone.

## Location

```go
Steps: []resource.TestStep{
    {
        ResourceName: "powerplatform_connection_shares.all_shares",
        Config: `
            data "azuread_domains" "aad_domains" {
                only_initial = true
            }
            // ... rest of the big HCL string
        `,
        Check: ...,
    },
}
```

## Code Issue

```go
Config: `
    data "azuread_domains" "aad_domains" {
        only_initial = true
    }
    ...
    resource "powerplatform_connection" "azure_openai_connection" {
        ...
        connection_parameters = jsonencode({
          "azureOpenAIResourceName" : "aaa",
          "azureOpenAIApiKey" : "bbb"
          "azureSearchEndpointUrl" : "ccc",
          "azureSearchApiKey" : "ddd"
        })
        ...
    }
    ...
`,
```

## Fix

Refactor the configuration blocks into clearly separated, smaller helper functions or constants. Store complex config as separate files or at least Go constants for easier editing and reuse. Optionally, use string formatting to minimize inlined logic and variable concatenations.

```go
const testAccConnectionShareConfig = `
data "azuread_domains" "aad_domains" {
    only_initial = true
}
...
`

Steps: []resource.TestStep{
    {
        ResourceName: "powerplatform_connection_shares.all_shares",
        Config: fmt.Sprintf(testAccConnectionShareConfig, mocks.TestsEntraLicesingGroupName(), mocks.TestName(), ...),
        Check: ...,
    },
}
```

Or, for long configs, maintain them in a directory as test fixtures and load them:

```go
configBytes, err := os.ReadFile("test-fixtures/connection_share_config.tf")
require.NoError(t, err)
...
Config: string(configBytes)
```
