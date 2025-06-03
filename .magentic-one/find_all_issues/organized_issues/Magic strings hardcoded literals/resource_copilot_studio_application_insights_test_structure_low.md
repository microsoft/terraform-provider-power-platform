# Title

Hardcoded Test Data Paths Reduce Maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights_test.go

## Problem

File system paths such as `"tests/Validate_Create/get_environments.json"` and `"tests/Test_Files/testagent_1_0_0_1_managed.zip"` are hardcoded as string literals in test configurations and HTTP mock file responses. This makes moving or restructuring the test data directories difficult and incurs a risk if paths need to be updated project-wide, as string search and replace may miss references.

## Impact

Low severity. This mainly impacts maintainability. When folder structures or test asset names change, the risk of stale test references increases.

## Location

Multiple test blocks, for example:

```go
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/get_environments.json").String()), nil
		})

	...

	resource "powerplatform_solution" "solution" {
		environment_id = powerplatform_environment.environment.id
		solution_file  = "tests/Test_Files/testagent_1_0_0_1_managed.zip"
	}
```

## Code Issue

```go
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/get_environments.json").String()), nil
		})

// and in test configuration strings:
solution_file  = "tests/Test_Files/testagent_1_0_0_1_managed.zip"
```

## Fix

Centralize test data paths as constants or use a function for test data path resolution. Example:

```go
const (
	testAgentZipPath           = "tests/Test_Files/testagent_1_0_0_1_managed.zip"
	validateCreateEnvsJsonPath = "tests/Validate_Create/get_environments.json"
)

...

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(validateCreateEnvsJsonPath).String()), nil
```

Or, for extra robustness, encapsulate in a test utility:

```go
func testDataPath(paths ...string) string {
	return filepath.Join(append([]string{"tests"}, paths...)...)
}
...
httpmock.File(testDataPath("Validate_Create", "get_environments.json"))
```

---

Continuing to scan for additional issues.