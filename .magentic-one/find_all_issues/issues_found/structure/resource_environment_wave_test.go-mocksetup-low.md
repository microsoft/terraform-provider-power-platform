# Repeated mock registration boilerplate reduces maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave_test.go

## Problem

Within the various test functions (Create, Error, NotFound, FailedDuringUpgrade, UnsupportedState), there is repeated, nearly identical code for registering environment, organizations, and feature endpoint mocks. Only the folder name and sometimes the mock data file differ. This violates DRY (Don't Repeat Yourself) principles and makes it easy to introduce inconsistencies or errors when editing or expanding the suite.

## Impact

Low. Readability, maintainability, and future extensibility decrease as the suite grows; subtle errors may be introduced if one helper registration is not correctly refactored everywhere.

## Location

Any test, for example:

```go
mocks.ActivateEnvironmentHttpMocks()
registerEnvironmentMock(t, "EnvironmentWaveResource_Create")

// Register organizations mock
registerOrganizationsMock(t, "EnvironmentWaveResource_Create")

// Register enable endpoint
httpmock.RegisterResponder("POST", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, ""), nil
	})
```

Similar code appears in each test.

## Fix

Factor repetitive mock setup logic into a shared helper function, parameterizing the test folder and mock filename(s):

```go
func setupEnvironmentWaveMocks(t *testing.T, testFolder string, featureResponder httpmock.Responder) {
	mocks.ActivateEnvironmentHttpMocks()
	registerEnvironmentMock(t, testFolder)
	registerOrganizationsMock(t, testFolder)
	httpmock.RegisterResponder(
		"POST",
		`=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features/October2024Update/enable$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		},
	)
	httpmock.RegisterResponder(
		"GET",
		`=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
		featureResponder,
	)
}
```

Each test would replace manual registration with a call like:

```go
setupEnvironmentWaveMocks(
	t,
	"EnvironmentWaveResource_FailedDuringUpgrade",
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_FailedDuringUpgrade", "get_features_failed.json")), nil
	})
```

This minimizes copy-paste and makes updates easier.
