# Title

Duplicated Test for Create Connection in Different Test Modes

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_test.go

## Problem

There are two test functions, `TestAccConnectionsResource_Validate_Create` and `TestUnitConnectionsResource_Validate_Create`, that effectively test nearly the same functional flow (creating a connection with specified parameters), only differing in their usage scenario: acceptance vs. unit test.

## Impact

While not strictly wrong, this duplication can make maintenance harder—if the connection creation flow changes, both tests must be updated. The test logic duplication could also be factored out to a helper for improved maintainability and readability.

Severity: Low

## Location

```go
func TestAccConnectionsResource_Validate_Create(t *testing.T) {
    ...
    // resource.Test(...)
    // Config and test steps...
}
func TestUnitConnectionsResource_Validate_Create(t *testing.T) {
    ...
    // httpmock
    // resource.Test(...)
    // Config and test steps...
}
```

## Fix

Factor out the shared test configuration fragment (HCL and check steps) into a helper function and use it in both tests.

```go
func connectionResourceConfig(name string) string {
    return fmt.Sprintf(`
        resource "powerplatform_connection" "azure_openai_connection" {
            environment_id = "%s"
            name           = "shared_azureopenai"
            display_name   = "OpenAI Connection %s"
            connection_parameters = jsonencode({
                "azureOpenAIResourceName" : "aaa",
                "azureOpenAIApiKey" : "bbb",
                "azureSearchEndpointUrl" : "ccc",
                "azureSearchApiKey" : "ddd"
            })
            lifecycle {
                ignore_changes = [
                    connection_parameters
                ]
            }
        }
    `, name, name)
}

// Use connectionResourceConfig(...) in both test steps as Config
```

This simplifies future changes and avoids copy-paste bugs.
