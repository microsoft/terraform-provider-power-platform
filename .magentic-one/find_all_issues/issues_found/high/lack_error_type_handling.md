# Issue: Lack of Error Type Handling in HTTP Mock Responses

### Path
`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go`

### Problem
HTTP mock functions do not provide structured error handling. All errors are handled generically, making debugging difficult if tests fail unexpectedly.

### Severity
High

### Suggested Fix
Expand mock responder behavior to include structured error handling and logging for improved transparency.

### Proposed Code Change
```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/...`, func(req *http.Request) (*http.Response, error) {
    response, err := httpmock.File(mockPath).String()
    if err != nil {
        log.Printf("Mock responder error: %v", err)
        return nil, fmt.Errorf("failed to load mocked response: %w", err)
    }
    return httpmock.NewStringResponse(http.StatusOK, response), nil
})
```
