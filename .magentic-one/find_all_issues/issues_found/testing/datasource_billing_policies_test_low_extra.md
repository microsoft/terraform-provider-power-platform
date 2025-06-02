# Title

No Assertions for External Provider Configuration Errors in Acceptance Test

## 

/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies_test.go

## Problem

In `TestAccBillingPoliciesDataSource_Validate_Read`, there is an implicit assumption that the `azapi` external provider will be available and the constraint will be satisfied. However, if the provider is not installed or the version constraint is violated, the test will fail in a non-obvious way with a generic error not directly related to the subject under test.

## Impact

Severity: **low**

This reduces test maintainability and may make failures harder to interpret, especially for new contributors or when running tests in new environments. Test failures appearing unrelated to the actual provider being tested can confuse diagnostics and slow down development.

## Location

First test function, near:

```go
ExternalProviders: map[string]resource.ExternalProvider{
	"azapi": {
		VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
		Source:            "azure/azapi",
	},
},
```

## Code Issue

```go
ExternalProviders: map[string]resource.ExternalProvider{
	"azapi": {
		VersionConstraint: constants.AZAPI_PROVIDER_VERSION_CONSTRAINT,
		Source:            "azure/azapi",
	},
},
```

## Fix

Add a test assertion or check at the beginning of the test to assert that the provider binary is present and satisfies the constraint, or explain this dependency in a test comment. Example with a comment for future maintainers:

```go
// NOTE: This test assumes the azapi external provider is present. If you see
// provider not found/version errors, ensure you have installed it per the
// README instructions or update the AZAPI_PROVIDER_VERSION_CONSTRAINT.
```
Alternatively, a pre-check function can assert provider prerequisites programmatically and fail fast with an informative error message.