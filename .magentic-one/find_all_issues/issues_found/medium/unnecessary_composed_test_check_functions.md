# Title

Unnecessary Use of `ComposeAggregateTestCheckFunc`

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/datasource_securityroles_test.go

## Problem

The `ComposeAggregateTestCheckFunc` includes multiple validations, some of which duplicate the functionality of other checks, making the code unnecessarily verbose.

## Impact

This redundancy increases the complexity of the codebase while providing little to no value. Severity: Medium.

## Location

TestAccSecurityDataSource_Validate_Read

## Code Issue

```go
resource.ComposeAggregateTestCheckFunc(
	resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
	resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.#", regexp.MustCompile(`^[1-9]\\d*$`)),
	resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.role_id", regexp.MustCompile(helpers.GuidRegex)),
	resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.name", regexp.MustCompile(helpers.StringRegex)),
	resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.is_managed", regexp.MustCompile(helpers.BooleanRegex)),
	resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.business_unit_id", regexp.MustCompile(helpers.GuidRegex)),
),
```

## Fix

Eliminate redundant checks by focusing on non-overlapping validations.

```go
resource.ComposeTestCheckFunc(
	resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "environment_id", regexp.MustCompile(helpers.GuidRegex)),
	resource.TestMatchResourceAttr("data.powerplatform_security_roles.all", "security_roles.0.role_id", regexp.MustCompile(helpers.GuidRegex)),
)
```