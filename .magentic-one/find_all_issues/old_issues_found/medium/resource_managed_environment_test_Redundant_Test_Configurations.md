# Redundant Test Configurations Across Steps

## Severity: Medium

### Location
File: `resource_managed_environment_test.go`
Function: `TestAccManagedEnvironmentsResource_Validate_Update`

### Description
In this test function, the same resource configuration is repeated with minor differences across multiple test steps. This results in redundant definitions and makes the test harder to maintain and extend.

### Impact
- Increases code duplication and reduces maintainability.
- Any changes to the base configuration require updates to multiple steps.

### Recommended Fix
Extract the common configuration into a reusable helper function or define a base configuration string that can be overridden with differences for specific test steps.

### Example of Fix
Define a helper function to generate the configuration:
```go
func generateTestConfig(baseConfig string, overrides map[string]string) string {
    config := baseConfig
    for key, value := range overrides {
        config = strings.Replace(config, fmt.Sprintf("{{%s}}", key), value, -1)
    }
    return config
}
```
Usage:
```go
baseConfig := `
    resource "example_resource" "test" {
        parameter_one = "{{parameter_one}}"
        parameter_two = "{{parameter_two}}"
    }
`
overrides := map[string]string{"parameter_one": "value1", "parameter_two": "value2"}
config := generateTestConfig(baseConfig, overrides)
```
