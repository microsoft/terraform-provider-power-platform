# Title

Hardcoded IDs in Resource Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install_test.go

## Problem

The configurations in the test files often use hardcoded IDs such as `00000000-0000-0000-0000-000000000001`:

```go
resource "powerplatform_environment_application_package_install" "development" {
	environment_id = "00000000-0000-0000-0000-000000000001"
	unique_name    = "ProcessMiningAnchor"
}
```

Hardcoding values in test cases reduces flexibility and overrides true dynamic testing environments with static dependencies.

## Impact

- **Critical Severity**: Hardcoded IDs can lead to unexpected test failures when migrating the test setup to another environment or when testing in a different scope.
- It actively limits the reusability of test scripts across domains due to the reliance on specific static values.

## Location

Example instance:

```go
httpmock.RegisterResponder("POST", "https://api.powerplatform.com/appmanagement/environments/00000000-0000-0000-0000-000000000001/applicationPackages/ProcessMiningAnchor/install?api-version=2022-03-01-preview",
	func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusAccepted, "")
		resp.Header.Add("Operation-Location", "https://api.powerplatform.com/appmanagement/environments/402c2b45-f5dc-e561-869f-368544f94a13/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1")
		return resp, nil
	})
```

Another instance:

```go
resource "powerplatform_environment_application_package_install" "development" {
	environment_id = "00000000-0000-0000-0000-000000000001"
	unique_name    = "MicrosoftFormsPro"
}
```

## Fix

Replace the hardcoded values with variables or dynamic mock configurations. For instance, utilize randomly generated test data:

```go
var generatedGUID = helpers.GenerateGUID()

resource "powerplatform_environment_application_package_install" "development" {
	environment_id = generatedGUID
	unique_name    = "ProcessMiningAnchor"
}
```

For HTTP mock responders, pass dynamically generated values using utility methods provided by `httpmock` or similar libraries:

```go
httpmock.RegisterResponder("POST", fmt.Sprintf("https://api.powerplatform.com/appmanagement/environments/%s/applicationPackages/ProcessMiningAnchor/install?api-version=2022-03-01-preview", generatedGUID),
	func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusAccepted, "")
		resp.Header.Add("Operation-Location", fmt.Sprintf("https://api.powerplatform.com/appmanagement/environments/%s/operations/475af49d-9bca-437f-8be1-9e467f44be8a?api-version=1", generatedGUID))
		return resp, nil
	})
```

This promotes modularity and ensures tests are more robust and adaptable across environments.