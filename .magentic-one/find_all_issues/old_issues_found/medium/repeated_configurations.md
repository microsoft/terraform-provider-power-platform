# Issue: Repetition of Configuration in Test Cases

### Path
`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go`

### Problem
Configuration blocks for `powerplatform_user` resources and dependencies are repeated across multiple test cases, violating the DRY principle.

### Severity
Medium

### Suggested Fix
Extract repeated configuration blocks into reusable functions or template files to improve code maintainability.

### Proposed Code Change
```go
func CreateUserResourceConfiguration() string {
    return fmt.Sprintf(`
    resource "azuread_user" "test_user" {
        user_principal_name = "` + mocks.TestName() + `@mockDomain.com"
        ...
    }
    resource "powerplatform_user" "new_user" {
        ...
    }
    `)
}
Config := CreateUserResourceConfiguration()
```
