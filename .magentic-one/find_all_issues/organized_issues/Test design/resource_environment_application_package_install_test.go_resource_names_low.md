# Unclear Variable and Resource Naming in Test Configurations (Naming)

##

/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install_test.go

## Problem

There are resource names like `"env"`, `"development"`, `"environment"`, `"env"`, etc., reused or inconsistently named across multiple test configurations, which makes reading and debugging failed tests confusing.

## Impact

Severity: Low

Unclear or inconsistent naming does not break function, but it reduces clarity and makes diagnosis harder for future maintainers.

## Location

In every test configuration block in the test steps.

## Code Issue

```hcl
resource "powerplatform_environment" "env" {
	display_name     = "displayname"
	location         = "europe"
	environment_type = "Sandbox"
}
resource "powerplatform_environment_application_package_install" "development" {
	environment_id   = powerplatform_environment.env.id
	unique_name      = "ProcessMiningAnchor"
}
```

## Fix

Use descriptive and unique resource names for each test scenario. For example, instead of "development" and "env", use "test_env_install" or "env_no_dataverse", etc.

```hcl
resource "powerplatform_environment" "test_install" {
    display_name     = "displayname"
    location         = "europe"
    environment_type = "Sandbox"
}
resource "powerplatform_environment_application_package_install" "test_package_install" {
    environment_id   = powerplatform_environment.test_install.id
    unique_name      = "ProcessMiningAnchor"
}
```
