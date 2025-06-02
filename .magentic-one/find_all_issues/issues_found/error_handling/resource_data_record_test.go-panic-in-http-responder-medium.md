# Panic Usage in Test HTTP Responder Instead of Proper Error Handling

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

In the HTTP responder used for simulating POST requests to the contacts endpoint (inside `TestUnitDataRecordResource_Validate_Update_Relationships`), there is direct usage of `panic` when handling errors while reading the request body.

```go
if err != nil {
    panic(err)
}
```

Using `panic` in a test handler is discouraged because it will cause the *entire* test suite (and maybe others running in the same process) to exit immediately, rather than providing a meaningful failure message at the level of the Go testing system. This can mask the root cause of failures, make debugging harder, and break test parallelism/CI runs unexpectedly.

## Impact

- **Abrupt test termination**: The whole test process may end abruptly, preventing graceful assertion/failure reporting.
- **Poor diagnostics**: Panics don't provide clear context compared to `t.Fatalf` or returning an error in test code.
- **Test unpredictability**: May cause partial output, test hangs or resource leaks in some CI environments.
- **Severity**: Medium

## Location

```go
bodyAsBytes := make([]byte, req.ContentLength)
_, err := req.Body.Read(bodyAsBytes)
if err != nil {
    panic(err)
}
bodyAsString := string(bodyAsBytes)
```

## Fix

In HTTP test handlers (and test code in general), avoid panics: report explicit test failures and/or return an appropriate HTTP 500 mock response with diagnostic output. 

For example:

```go
import "io"
import "net/http"

httpmock.RegisterResponder("POST", "...",
    func(req *http.Request) (*http.Response, error) {
        bodyAsBytes, err := io.ReadAll(req.Body)
        if err != nil {
            // Return a 500 error with diagnostic content instead of panicking
            return httpmock.NewStringResponse(http.StatusInternalServerError, err.Error()), nil
        }
        bodyAsString := string(bodyAsBytes)
        // ... switch ...
    })
```

Or, if you want explicit test failure notification, you could inject a `*testing.T` parameter and call `t.Fatalf`, but avoid raw panics wherever possible.

---
