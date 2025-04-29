# Issue: Excessive Test Dependencies on External Resources

### Path
`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go`

### Problem
Tests rely heavily on external dependencies such as AzureAD Terraform providers. These dependencies can introduce flakiness into tests in case the external services or providers are unavailable or unstable.

### Severity
Critical

### Suggested Fix
- Replace external resources with mock implementations or abstractions that mimic their behavior without depending on actual providers.

### Proposed Code Change
```go
Config: `
mock_data "azuread_domains" "aad_domains" {
    domains = ["mockDomain.com"]
}
resource "mock_azuread_user" "test_user" {
    user_principal_name = "` + mocks.TestName() + `@mockDomain.com"
}
`
```
