# Hardcoded Mock Endpoint URLs in Tests

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set_test.go

## Problem

The test file contains many occurrences of hardcoded API URLs in the `httpmock.RegisterResponder` setup (for example: `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/...`). These endpoints are repeated across different test functions, and the GUIDs are fixed, which may lead to maintenance overhead if endpoints or tenants/GUIDs change. Furthermore, these hardcoded test values can create confusing debugging error messages or limit reusability if more dynamic or parametrized test coverage is required.

## Impact

Severity: Low

While not an immediate logic bug, this reduces test maintainability, increases duplication, and may lead to missed updates whenever endpoint structures change in the provider.

## Location

All `httpmock.RegisterResponder` calls (e.g., lines around "POST", "GET", "PUT", "DELETE" registration):

```go
httpmock.RegisterResponder("POST", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000000/ruleSets?api-version=2021-10-01-preview`,
...
```

## Fix

Refactor repeated endpoint URLs into constants, and/or provide helpers for building these URLs. This helps reduce duplication and makes bulk changes simple and less error-prone.

```go
const (
	testTenantID     = "000000000000000000000000000000.01"
	testGroupID      = "00000000-0000-0000-0000-000000000000"
	testRuleSetID    = "00000000-0000-0000-0000-000000000001"
	apiVersionSuffix = "?api-version=2021-10-01-preview"
)

var (
	baseAPI = fmt.Sprintf("https://%s.tenant.api.powerplatform.com/governance", testTenantID)
	getRuleSetsURL = fmt.Sprintf("%s/environmentGroups/%s/ruleSets%s", baseAPI, testGroupID, apiVersionSuffix)
	...
)

httpmock.RegisterResponder("POST", getRuleSetsURL, ...)
```

Add a comment explaining that test GUIDs and tenants are intentionally static for test isolation.
