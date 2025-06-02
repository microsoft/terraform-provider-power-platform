# Title

Incomplete Error Handling for HTTPMock

## Path

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

In cases where HTTPMock responders are registered, there is no explicit error handling to ensure HTTP responders behave correctly or to catch runtime problems that might arise during mock setup.

For example, within certain responder callbacks like:
```go
httpmock.RegisterResponder("PATCH", fmt.Sprintf(`https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/%s`, policyId),
    func(req *http.Request) (*http.Response, error) {
        patchResponsesInx++
        return httpmock.NewStringResponse(http.StatusOK, patchResponsesArray[patchResponsesInx]), nil
    })
```

Failures such as an `index out-of-range` for `patchResponsesArray` will lead to runtime panics instead of being silently handled.

## Impact

Severity is **critical** as unchecked assumptions can cause runtime crashes or incorrect behaviors in unit tests, interrupting CI/CD workflows, and misrepresenting test reliability.

## Location

Found in responders using callbacks such as:
- `TestUnitDataLossPreventionPolicyResource_Validate_Update`
- `TestUnitDataLossPreventionPolicyResource_Validate_Create`

## Code Issue

```go
patchResponsesInx++
return httpmock.NewStringResponse(http.StatusOK, patchResponsesArray[patchResponsesInx]), nil
```

## Fix

Introduce error handling or bounds checking around array access and validate input state thoroughly before incrementing indexes.

### Example Fix

```go
httpmock.RegisterResponder("PATCH", fmt.Sprintf(`https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies/%s`, policyId),
    func(req *http.Request) (*http.Response, error) {
        if patchResponsesInx+1 >= len(patchResponsesArray) {
            return nil, fmt.Errorf("index out of range for patchResponsesArray")
        }
        patchResponsesInx++
        return httpmock.NewStringResponse(http.StatusOK, patchResponsesArray[patchResponsesInx]), nil
    })
```