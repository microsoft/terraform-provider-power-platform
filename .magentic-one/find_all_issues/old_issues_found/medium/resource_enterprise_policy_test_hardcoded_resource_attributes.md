# Title

Hardcoded Resource Attributes in Test Configuration

##

`/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy_test.go`

## Problem

The test case includes hardcoded resource attributes (e.g., `display_name`, `description`, `location`, etc.). This practice limits the reusability of the test case and makes it hard to scale or adapt to different use cases.

## Impact

Hardcoded values reduce test flexibility and may lead to duplication when additional test cases requiring slight variations need to be added. This hampers maintainability. Severity: **Medium**

## Location

The hardcoded configuration is in the `Config` parameter of the `resource.TestStep`:

```go
Config: `
	resource "powerplatform_environment" "development" {
		display_name                              = "displayname"
		description                               = "description"
		cadence									      = "Frequent"
		location                                  = "europe"
		environment_type                          = "Sandbox"
	}

	resource "powerplatform_enterprise_policy" "network_injection" {
		environment_id = powerplatform_environment.development.id
		system_id      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"
		policy_type    = "NetworkInjection"
	}`,
```

---

## Code Issue

```go
resource "powerplatform_environment" "development" {
	display_name                              = "displayname"
	description                               = "description"
	cadence									      = "Frequent"
	location                                  = "europe"
	environment_type                          = "Sandbox"
}

resource "powerplatform_enterprise_policy" "network_injection" {
	environment_id = powerplatform_environment.development.id
	system_id      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"
	policy_type    = "NetworkInjection"
}
```

## Fix

Refactor the test configuration to use variables or constants at the top of the test file for the resource attributes. These can then be dynamically passed into the configuration string, making the test more flexible and easier to maintain.

```go
const displayName = "displayname"
const description = "description"
const cadence = "Frequent"
const location = "europe"
const environmentType = "Sandbox"
const systemID = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"

Config: fmt.Sprintf(`
	resource "powerplatform_environment" "development" {
		display_name                              = "%s"
		description                               = "%s"
		cadence									      = "%s"
		location                                  = "%s"
		environment_type                          = "%s"
	}

	resource "powerplatform_enterprise_policy" "network_injection" {
		environment_id = powerplatform_environment.development.id
		system_id      = "%s"
		policy_type    = "NetworkInjection"
	}`, displayName, description, cadence, location, environmentType, systemID)
```

This refactor enhances test flexibility and maintainability.
