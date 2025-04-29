# Issue: Hardcoded Mock Paths in Test Code

### Path
`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go`

### Problem
Several tests rely on hardcoded file paths for mocking responses, such as `httpmock.File("tests/resource/user/...")`. These hardcoded paths may fail in environments where the relative path configuration is different, leading to brittle tests.

### Severity
Medium

### Suggested Fix
- Abstract mock file paths into test configuration variables or constants.
- Use a standardized function to resolve file paths based on the current working directory or test-specific configuration.

### Proposed Code Change
```go
// Replace hardcoded path with a configurable/test-agnostic path resolver function
mockPath := ResolveMockFilePath("tests/resource/user/...")
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/...`, func(req *http.Request) (*http.Response, error) {
    return httpmock.NewStringResponse(http.StatusOK, httpmock.File(mockPath).String()), nil
})
```
