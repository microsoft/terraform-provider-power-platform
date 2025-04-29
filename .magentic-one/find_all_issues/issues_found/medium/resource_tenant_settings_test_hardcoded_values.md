# Title

Hard-Coded Values in Test Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings_test.go

## Problem

The test configurations contain multiple hard-coded values for tenant settings attributes. Using hard-coded values in tests makes it difficult to adapt the tests to changes in the application or environment, leading to potential brittleness.

## Impact

- Reduces flexibility: Changes in the configuration may require multiple modifications in the test scripts.
- Makes tests less reusable: Cannot parameterize easily for other use cases or environments.
- Medium severity: While this does not completely break functionality, it hampers test maintainability.

## Location

Throughout `TestAccTenantSettingsResource_Validate_Create` and `TestAccTenantSettingsResource_Validate_Update` functions.

## Code Issue

Hard-coded value example:

```go
disable_support_tickets_visible_by_all_users = false
disable_trial_environment_creation_by_non_admin_users = false
disable_capacity_allocation_by_environment_admins = true
disable_environment_creation_by_non_admin_users = false
```

## Fix

Replace the hard-coded values with variables that can be parameterized. This allows developers to modify settings for different test cases more easily.

```go
var testConfig = map[string]string{
	"disable_support_tickets_visible_by_all_users": "false",
	"disable_trial_environment_creation_by_non_admin_users": "false",
	"disable_capacity_allocation_by_environment_admins": "true",
	"disable_environment_creation_by_non_admin_users": "false",
}

func getTestConfig(config map[string]string) string {
	// Construct configuration dynamically using provided parameters.
	return fmt.Sprintf(`
	resource "powerplatform_tenant_settings" "settings" {
		disable_support_tickets_visible_by_all_users = %s
		disable_trial_environment_creation_by_non_admin_users = %s
		disable_capacity_allocation_by_environment_admins = %s
		disable_environment_creation_by_non_admin_users = %s
	}
	`, config["disable_support_tickets_visible_by_all_users"], 
		config["disable_trial_environment_creation_by_non_admin_users"], 
		config["disable_capacity_allocation_by_environment_admins"], 
		config["disable_environment_creation_by_non_admin_users"])
}
```