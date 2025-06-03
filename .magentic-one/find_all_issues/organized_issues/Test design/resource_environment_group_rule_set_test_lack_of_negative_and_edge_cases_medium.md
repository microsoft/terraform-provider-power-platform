# Lack of Negative and Edge Case Testing

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set_test.go

## Problem

The test cases, both in acceptance (when enabled) and unit layers, only cover "happy path" CRUD scenarios with valid data. They do not test negative paths such as malformed configurations, API rejections, required-field omission, concurrency, or permission failures. This reduces confidence in resilience and validation, especially as API contracts or requirements might change.

## Impact

Severity: Medium

A lack of negative and edge-case tests can allow subtle bugs (e.g., bad input handling, missing field errors, incorrect error messaging) to escape to production. This reduces overall confidence in implementation correctness.

## Location

All acceptance and unit tests are affected; examples:

```go
resource.TestCheckResourceAttr("...", ...)
// but no test of error conditions, only valid configs.
```

## Fix

Add tests for invalid or missing fields, API failure scenarios, and field range/constraint violations.

Example:

```go
{
	Config: `
		resource "powerplatform_environment_group_rule_set" "example" {
			// missing environment_group_id (required)
			rules = {...}
		}
	`,
	ExpectError: regexp.MustCompile(".*environment_group_id.*required.*"),
},
{
	Config: `
		resource "powerplatform_environment_group_rule_set" "example" {
			environment_group_id = "bad-guid"
			// ...
		}
	`,
	ExpectError: regexp.MustCompile(".*must be a valid UUID.*"),
},
```

Create additional `httpmock.RegisterResponder` calls that simulate error HTTP statuses and verify correct error handling.
