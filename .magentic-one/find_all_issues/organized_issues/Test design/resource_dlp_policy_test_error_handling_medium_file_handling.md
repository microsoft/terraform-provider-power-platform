# Potentially Unchecked Errors from httpmock.File().String()

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

In the test, the code uses `httpmock.File(...).String()` directly as the body for a mock HTTP response. If the file does not exist or cannot be read, this could panic or lead to misleading test errors. There is no error handling to capture file-not-found or read errors, which can obscure underlying file system/test fixture issues.

## Impact

**Severity: Medium**  
Uncaught errors in reading files can cause tests to panic unexpectedly, making debugging harder and hiding the root cause of test failures. It also makes the test fragile if test fixture files are renamed, relocated, or corrupted.

## Location

Inside HTTP responder setup for POST/GET:

```go
httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v2/policies`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/Validate_Create/get_policy_00000000-0000-0000-0000-000000000001.json").String()), nil
    })
```

## Fix

Capture and explicitly handle possible errors from reading the file (if supported by `httpmock.File`—otherwise, improve testability by using helper utilities that report or fail with an informative message):

```go
func readFileOrFail(t *testing.T, filename string) string {
    data, err := os.ReadFile(filename)
    if err != nil {
        t.Fatalf("failed to read test fixture: %v", err)
    }
    return string(data)
}

// Then in responder:
func(req *http.Request) (*http.Response, error) {
    return httpmock.NewStringResponse(http.StatusCreated, readFileOrFail(t, "tests/resource/Validate_Create/get_policy_...json")), nil
}
```

If `httpmock.File` cannot return errors, this is a limitation and should be documented, or a wrapper should be created to avoid surprise test panics.
