# Magic Strings and Inlined Configuration in Test Code

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query_test.go

## Problem

Provider configuration for Terraform is written as inline, multi-line strings. Magic values for things like IDs, names, and API URLs are hard-coded, making it harder to refactor or reuse across multiple tests.

## Impact

Medium severity for maintainability. Changes in example values require search-and-replace rather than a global update. Greater risk of copy-paste errors.

## Location

```go
Config: `
	resource "powerplatform_environment" "env" {
		display_name     = "` + mocks.TestName() + `"
		location         = "unitedstates"
		environment_type = "Sandbox"
		dataverse = {
			language_code     = "1033"
			currency_code     = "USD"
			security_group_id = "00000000-0000-0000-0000-000000000000"
		}
	}

	data "powerplatform_rest_query" "webapi_query" {
		scope                = "${powerplatform_environment.env.dataverse.url}/.default"
		url                  = "${powerplatform_environment.env.dataverse.url}api/data/v9.2/WhoAmI"
		method               = "GET"
		expected_http_status = [200]
	}
`,
```

## Fix

Extract reusable parts, IDs, and sample names to constants:

```go
const (
	testEnvironmentName   = "TestEnvironment"
	testLocation          = "unitedstates"
	testEnvironmentType   = "Sandbox"
	testLanguageCode      = "1033"
	testCurrencyCode      = "USD"
	testSecurityGroupID   = "00000000-0000-0000-0000-000000000000"
)

Config: fmt.Sprintf(`
	resource "powerplatform_environment" "env" {
		display_name     = "%s"
		location         = "%s"
		environment_type = "%s"
		dataverse = {
			language_code     = "%s"
			currency_code     = "%s"
			security_group_id = "%s"
		}
	}
	// ...
`, testEnvironmentName, testLocation, testEnvironmentType, testLanguageCode, testCurrencyCode, testSecurityGroupID)
```

This improves maintainability and readability.
