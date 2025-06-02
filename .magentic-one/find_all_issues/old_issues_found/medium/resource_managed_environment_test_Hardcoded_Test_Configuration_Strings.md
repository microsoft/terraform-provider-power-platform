# Hardcoded Test Configuration Strings in Resource Definition

## Severity: Medium

### Location
File: `resource_managed_environment_test.go`
Function: `TestAccManagedEnvironmentsResource_Validate_Create`

### Description
Configuration strings that define the test resource setup are hardcoded directly inline within the test function using a multiline string format. This approach impairs test readability, increases the likelihood of errors, and makes it harder to maintain consistency across multiple test cases that might re-use the same resource configurations.

### Impact
- Reusability is compromised, and duplication of configurations across tests can lead to inconsistent setups.
- If the configuration format changes, it requires manual updates wherever it is reused.

### Recommended Fix
Define reusable configuration structures or utilize a helper function to generate the configuration dynamically. This approach improves maintainability and ensures consistency.

### Examples of Fix
Create a helper function named `generate_test_config`:
```go
func generate_test_config(name string, value string) string {
    return fmt.Sprintf(`
        resource "example_resource" "%s" {
            parameter_one = "%s"
        }
    `, name, value)
}
```
Usage:
```go
config := generate_test_config("test_example", "test_value")
```
