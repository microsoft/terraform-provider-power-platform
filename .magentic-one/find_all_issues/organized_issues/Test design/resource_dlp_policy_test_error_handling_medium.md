# Use of Global Mutable State for Test Response Indexing

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

The variables `getResponsesInx` and `patchResponsesInx` are used as global counters and mutated by HTTP responder functions to determine which mock response to return. This approach introduces hidden state and can cause unpredictable test behavior, especially if tests are run in parallel or reentrant, or if a panic prevents index reset.

## Impact

**Severity: Medium**  
This could lead to data races or incorrect responses when running tests concurrently or if more than one test ever relies on these responders. It also makes the test harder to reason about due to hidden mutable state.

## Location

Definition and usage at lines where the indexers are incremented and at their declaration:

```go
getResponsesInx := -1
// ...
patchResponsesInx := -1
// ...
func(req *http.Request) (*http.Response, error) {
	getResponsesInx++
	return httpmock.NewStringResponse(http.StatusOK, getResponsesArray[getResponsesInx]), nil
}
```

## Fix

Encapsulate response arrays and indices in a struct or closure to avoid global mutable state. Alternatively, use atomic operations or `t.Parallel()` avoidance to ensure safety, but the cleanest approach is closure-based.

```go
// Example: Closure-based indexer for responders
func indexedResponder(responses []string) func(req *http.Request) (*http.Response, error) {
	index := -1
	return func(req *http.Request) (*http.Response, error) {
		index++
		return httpmock.NewStringResponse(http.StatusOK, responses[index]), nil
	}
}

// ... Register as:
httpmock.RegisterResponder("GET", someURL, indexedResponder(getResponsesArray))
httpmock.RegisterResponder("PATCH", otherURL, indexedResponder(patchResponsesArray))
```
