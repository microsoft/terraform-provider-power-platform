# Title

Test does not assert HTTP mock usage or reset state

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share_test.go

## Problem

Although httpmock is used and reset with `defer httpmock.DeactivateAndReset()`, there are no assertions to verify that the registered HTTP mocks are called (and in the order/quantity expected). The test could pass even if the provider does not make the mocked calls due to a regression, because there is no assertion on HTTP mock usage.

## Impact

Tests may result in false positives—such as passing even when no HTTP communication occurs—resulting in a lack of full confidence in the correctness of the integration points. This issue is medium severity because it affects the reliability of test results.

## Location

Relevant to the `TestUnitConnectionsShareResource_Validate_Create` function.

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
// ... responders registered, but no assert httpmock.CallCountInfo()
resource.Test(t, resource.TestCase{ ... })
```

## Fix

After running the test, check the mock invocations using `httpmock.GetCallCountInfo()`. Assert that expected requests have been made:

```go
func TestUnitConnectionsShareResource_Validate_Create(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    
    // (register responders)

    resource.Test(t, resource.TestCase{ ... })

    // Assert HTTP mock was used as expected
    info := httpmock.GetCallCountInfo()
    expected := map[string]int{
        "POST https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/modifyPermissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1": 1,
        "GET https://000000000000000000000000000000.00.environment.api.powerplatform.com/connectivity/connectors/shared_commondataserviceforapps/connections/00000000-0000-0000-0000-000000000001/permissions?%24filter=environment+eq+%2700000000-0000-0000-0000-000000000000%27&api-version=1": 1,
    }

    for endpoint, count := range expected {
        if info[endpoint] != count {
            t.Errorf("Expected %d calls to %s, got %d", count, endpoint, info[endpoint])
        }
    }
}
```

This guarantees the resource uses the underlying API as intended during tests.
