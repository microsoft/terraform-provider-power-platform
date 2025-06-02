# Inefficient/Unsafe HTTP Mock Body Reading in Test

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

In the test function `TestUnitDataRecordResource_Validate_Update_Relationships`, a POST responder for the contacts endpoint reads the request body with:

```go
bodyAsBytes := make([]byte, req.ContentLength)
_, err := req.Body.Read(bodyAsBytes)
if err != nil {
    panic(err)
}
bodyAsString := string(bodyAsBytes)
```

This approach is unsafe and could easily produce partial reads. Go's `io.Reader` implementations (including `http.Request.Body`) are not guaranteed to fill the buffer in one call to `Read`—it often requires multiple reads. Using `Read` like this may result in an incomplete body under high concurrency, leading to subtle, non-deterministic bugs in your tests.

## Impact

- **Test Flakiness**: Test could sometimes work, sometimes fail, due to partial body reads (non-deterministic flakiness).
- **Incorrect Logic**: Mismatched body, broken test-case routing, panic if ContentLength isn't set.
- **Severity**: Medium (Risk of hiding regressions or enabling false positives in CI).

## Location

```go
httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts`,
    func(req *http.Request) (*http.Response, error) {
        bodyAsBytes := make([]byte, req.ContentLength)
        _, err := req.Body.Read(bodyAsBytes)
        if err != nil {
            panic(err)
        }
        bodyAsString := string(bodyAsBytes)

        // ... switch on bodyAsString ...
    })
```

## Fix

Always read the full body by using `io.ReadAll` (Go 1.16+), which is both correct and idiomatic:

```go
import "io"

// ...

httpmock.RegisterResponder("POST", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts`,
    func(req *http.Request) (*http.Response, error) {
        bodyAsBytes, err := io.ReadAll(req.Body)
        if err != nil {
            panic(err) // or return error if appropriate
        }
        bodyAsString := string(bodyAsBytes)

        // ... switch on bodyAsString ...
    })
```

If you need to reset or reuse the body, wrap with `io.NopCloser(bytes.NewBuffer(bodyAsBytes))`.

---
