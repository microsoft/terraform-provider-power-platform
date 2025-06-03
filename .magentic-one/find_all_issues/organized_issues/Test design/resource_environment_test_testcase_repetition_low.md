# Test Function Repetition and Bulkiness

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

There is significant repetition in test case setups and test scaffolding within this file. Many test functions contain duplicated patterns for activating mocks, registering HTTP responders, and constructing large config blocks with mostly identical content differing only by a few values. Large inlined configuration blocks (Terraform test configs) are repeated with only minor variations, making the test difficult to maintain and reason about.

## Impact

- **Severity: Low (for correctness), Medium (for maintainability)**
- Makes the file harder to maintain and increases the risk of copy-paste mistakes.
- Future changes (such as a config schema update or change in environment setup) would have to be repeated throughout many locations, raising the likelihood of inconsistencies and bugs.
- Test readability and navigability is reduced as patterns are not DRY (Don't Repeat Yourself) and repeated logic is not abstracted.

## Location

Throughout the file, for example:
- Repeated registration of the same HTTP mock endpoints in the `TestUnitEnvironmentsResource_Validate_*` functions
- Similar Terraform config snippets repeated in multiple test steps

## Code Issue

```go
resource.Test(t, resource.TestCase{
    // ...
    Steps: []resource.TestStep{
        {
            Config: ` // repeated config block
                resource "powerplatform_environment" ...
            `,
            Check: resource.ComposeTestCheckFunc(...),
        },
        // ...additional similar steps
    },
})
// ...
```

## Fix

Abstract common setup and configuration patterns using helper functions. Store repeated configuration templates in a variable, and use Go functions for instantiation. For mocks, use helper functions to register standard HTTP responders, and use parameterization for common config.

```go
// Example: Helper to provide standard config
func environmentConfig(displayName, location, envType string) string {
    return fmt.Sprintf(`
        resource "powerplatform_environment" "development" {
            display_name = "%s"
            location = "%s"
            environment_type = "%s"
            // ...
        }
    `, displayName, location, envType)
}

func registerStandardEnvironmentMocks() {
    httpmock.RegisterResponder(...) // etc
}

// Usage
registerStandardEnvironmentMocks()
config := environmentConfig("displayname", "europe", "Sandbox")
resource.Test(t, resource.TestCase{
    Steps: []resource.TestStep{{Config: config, Check: ...}},
})
```

This reduces file size, improves maintainability, and encourages consistency across tests.
