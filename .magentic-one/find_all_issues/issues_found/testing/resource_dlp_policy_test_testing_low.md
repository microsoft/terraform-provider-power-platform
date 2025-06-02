# Test Skipped Indefinitely Without Proper Issue Tracking or Reference

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

The acceptance test `TestAccDataLossPreventionPolicyResource_Validate_Create` is permanently skipped due to "inconsistency in API in connectors returned" but does not reference any external tracking (issue, ticket, TODO, etc.), nor does it provide details for future maintainers on what needs to be fixed or how/when to re-enable the test.

## Impact

**Severity: Low**  
This makes it likely that the test is ignored indefinitely and never re-enabled or properly investigated. It hampers test coverage and traceability of broken or unimplemented features.

## Location

The skip statement:

```go
func TestAccDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
    t.Skip("Skipping as there is inconsistency in API in connectors returned")

    resource.Test(t, resource.TestCase{
        // ...
    })
}
```

## Fix

Provide a reference to a GitHub issue or internal tracking number, or at least a TODO so future maintainers can track the status and motivation for skipping. Optionally, add more details explaining the issue and circumstances that would make the test re-enablable.

```go
func TestAccDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
    t.Skip("Skipping as there is inconsistency in API in connectors returned. " +
        "TODO: See issue #123 (or similar tracking reference) for resolution tracking.")

    // resource.Test(...) ...
}
```
