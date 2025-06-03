# Title

Use of Deprecated net/httpmock.File for Reading Test Data

## 

/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies_test.go

## Problem

In the unit test, the code uses `httpmock.File("test/datasource/policies/get_billing_policies.json")` to read the response body from a file when constructing an HTTP response. Reading files directly like this during tests is fragile and can cause test failures if the file is missing, renamed, or if the working directory changes. Furthermore, some http mocking libraries discourage or deprecate this usage in favor of loading fixtures with Go's `os.ReadFile` or embedding test data with Go 1.16+ `embed` directive.

## Impact

Severity: **low**

Direct file reads can cause non-deterministic test failures. Tests become dependent on correct relative file paths and project structure. If the file is missing or renamed, the test will panic or fail in a non-obvious way.

## Location

Second test function, HTTP responder construction:

## Code Issue

```go
httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/datasource/policies/get_billing_policies.json").String()), nil
	})
```

## Fix

Load the file using Go's standard library, preferably at package initialization time or test setup, and pass its contents explicitly. This can also be done via Go's `embed` for static fixtures:

**Using os.ReadFile in your test:**
```go
import (
	"os"
	...
)

data, err := os.ReadFile("test/datasource/policies/get_billing_policies.json")
if err != nil {
	t.Fatalf("failed to read billing policies fixture: %v", err)
}
httpmock.RegisterResponder("GET", "...",
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, string(data)), nil
	})
```
**Or using Go embed (1.16+):**
```go
import _ "embed"

//go:embed test/datasource/policies/get_billing_policies.json
var policyFixture string

httpmock.RegisterResponder("GET", "...",
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, policyFixture), nil
	})
```
This ensures your tests are less brittle and you get better compile-time or test startup-time errors if the file goes missing.
