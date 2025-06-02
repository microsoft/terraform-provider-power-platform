# Duplicate httpmock.RegisterResponder registrations reduce test clarity

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave_test.go

## Problem

Multiple `httpmock.RegisterResponder` calls are registered for the same HTTP method and URL pattern, particularly for several GET requests to the features endpoint. In Go's `httpmock`, the last registration for a given method/URL pattern "wins", shadowing previous responders for that route. This makes the intended state machine for simulating state transitions unclear—it's not explicit how to simulate sequential state changes (e.g., "Upgrading" to "On" state). The tests assume (undocumented) responder stack behavior, which may lead to confusion or incorrect test assumptions.

## Impact

Medium. Lack of explicitness in mock responses for sequence-based interactions makes the tests difficult to understand or maintain and can easily cause regressions if `httpmock` changes registration mechanics. Debugging is harder if mock state transitions do not behave as expected.

## Location

```go
// Example from TestUnitEnvironmentWaveResource_Create
httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_Create", "get_features_upgrading.json")), nil
	})

httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, loadTestResponse(t, "EnvironmentWaveResource_Create", "get_features_enabled.json")), nil
	})
```

## Fix

Use an `httpmock.SeqResponder` or similar pattern to simulate state transitions in a controlled, explicit way:

```go
var featureStates = []string{
	loadTestResponse(t, "EnvironmentWaveResource_Create", "get_features_upgrading.json"),
	loadTestResponse(t, "EnvironmentWaveResource_Create", "get_features_enabled.json"),
}
var idx int

httpmock.RegisterResponder("GET", `=~^https://api\.admin\.powerplatform\.microsoft\.com/api/environments/00000000-0000-0000-0000-000000000001/features$`,
	func(req *http.Request) (*http.Response, error) {
		// Serve the current state, advance for next call
		if idx >= len(featureStates) {
			idx = len(featureStates) - 1
		}
		resp := featureStates[idx]
		idx++
		return httpmock.NewStringResponse(http.StatusOK, resp), nil
	})
```

This clarification makes the test state progression clear and intentional, enhancing readability and maintainability.
