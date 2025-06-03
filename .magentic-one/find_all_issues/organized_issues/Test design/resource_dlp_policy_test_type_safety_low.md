# Lack of Proper Type Safety Between Test Terraform Config Strings and Expected Values

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

Test cases embed multi-line Terraform HCL config strings and assert resource attributes using string-based keys and expected values. While typical in Terraform-provider test code, this string-based convention is error-prone: typos in attribute names, wrong indices, or ordering mismatches will lead to misleading test failures—especially given the heavy use of complex nested and set/array data.

## Impact

**Severity: Low**  
This weakens the reliability and maintainability of the test suite, especially as schemas evolve. There’s no compile-time checking of attribute keys or their expected types, so failures may go unnoticed until runtime.

## Location

Throughout all the test cases:

```go
resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "non_business_connectors.0.id", "/providers/Microsoft.PowerApps/apis/shared_sql"),
resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", "custom_connectors_patterns.0.host_url_pattern", "https://*.contoso.com"),
```

and so on, with hand-coded indices and keys.

## Fix

While Terraform provider testing relies on this pattern, improve type safety by:

- Defining constants for attribute keys to avoid typos.
- For more complex validation, use custom check functions that parse and assert important struct fields directly, e.g., decode into structs and check properties rather than asserting by brittle attribute strings.
- Consider adding helper functions/macros (if possible in Go) to build and check configurations.

```go
const (
    attrNonBusinessConnectorsID           = "non_business_connectors.0.id"
    attrCustomConnectorsPatternsHostURL   = "custom_connectors_patterns.0.host_url_pattern"
    // ...
)

resource.TestCheckResourceAttr("powerplatform_data_loss_prevention_policy.my_policy", attrNonBusinessConnectorsID, "/providers/Microsoft.PowerApps/apis/shared_sql"),
```

For critical/complex objects, decode the resource state and use Go struct checks for type safety (requires more involved plumbing but provides strong assurances).
