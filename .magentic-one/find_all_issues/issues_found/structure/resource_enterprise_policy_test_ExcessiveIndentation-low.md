# Issue: Excessive Line Indentation in Test Configuration

## 
/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy_test.go

## Problem

The Terraform config in the test includes inconsistent tabbing and excessive indentation. This reduces readability and maintainability.

## Impact

Affects readability for those updating or reviewing test configurations. Severity: low.

## Location

Lines 61-82, in the `Config:` string.

## Code Issue

```go
Config: `
	resource "powerplatform_environment" "development" {
		display_name                              = "displayname"
		description                               = "description"
		cadence				      = "Frequent"
		location                                  = "europe"
		environment_type                          = "Sandbox"
	}

	resource "powerplatform_enterprise_policy" "network_injection" {
		environment_id = powerplatform_environment.development.id
		system_id      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"
		policy_type    = "NetworkInjection"
	}`,
```

## Fix

Use consistent 2- or 4-space indentation and replace tabs for alignment:

```go
Config: `
  resource "powerplatform_environment" "development" {
    display_name    = "displayname"
    description     = "description"
    cadence         = "Frequent"
    location        = "europe"
    environment_type = "Sandbox"
  }

  resource "powerplatform_enterprise_policy" "network_injection" {
    environment_id = powerplatform_environment.development.id
    system_id      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000002"
    policy_type    = "NetworkInjection"
  }
`,
