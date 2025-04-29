# Title

Improper Cleanup in Unit Test

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings_test.go

## Problem

The unit tests using the HTTP Mock library (`TestUnitTestTenantSettingsResource_Validate_Create` and `TestUnitTestTenantSettingsResource_Validate_Update`) do not clean up or validate the HTTP mock's state after tests are run. Although `httpmock.DeactivateAndReset` is called, there is no validation to ensure all registered responders are triggered correctly or removed.

## Impact

- Potential cross-contamination between tests, as mocks may carry over unexpected behaviors/results.
- Reduced reliability in test scenarios using mocks.
- Medium severity: This can cause flaky tests and debugging challenges during failures.

## Location

Mainly in `TestUnitTestTenantSettingsResource_Validate_Create` and `TestUnitTestTenantSettingsResource_Validate_Update`.

## Code Issue

Improper cleanup:

```go
httpmock.DeactivateAndReset() // Only deactivates without validating registered responders or resetting them correctly
```

## Fix

Validate that all mock responders are triggered and remove untriggered responders before resetting the mock state.

```go
func TestUnitTestTenantSettingsResource_Validate_Create(t *testing.T) {
	defer func() {
		httpmock.DeactivateAndReset()
		if len(httpmock.GetCallCountInfo()) > 0 {
			t.Errorf("Detected unused HTTP mocks: %v", httpmock.GetCallCountInfo())
		}
	}()

	httpmock.Activate()
	// Setup mock responders and test logic
}

func TestUnitTestTenantSettingsResource_Validate_Update(t *testing.T) {
	defer func() {
		httpmock.DeactivateAndReset()
		if len(httpmock.GetCallCountInfo()) > 0 {
			t.Errorf("Detected unused HTTP mocks: %v", httpmock.GetCallCountInfo())
		}
	}()

	httpmock.Activate()
	// Setup mock responders and test logic
}
```