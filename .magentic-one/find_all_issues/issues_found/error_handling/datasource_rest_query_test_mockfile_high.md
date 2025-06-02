# Error Handling in HTTP Mock Callback

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query_test.go

## Problem

The callback functions for `httpmock.RegisterResponder` directly access files with `httpmock.File(...).String()` without error handling. If the mock file does not exist or is unreadable, this will panic or pass an empty/invalid string, rather than failing with a helpful error.

## Impact

**High** severity in terms of testing reliability. A missing or corrupt fixture might silently make the test pass with wrong data or fail in a non-clear way.

## Location

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
```

## Code Issue

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
```

## Fix

Add explicit error handling for file reads:

```go
fileContents, err := os.ReadFile("tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json")
if err != nil {
	t.Fatalf("Failed to open mock file: %v", err)
}
return httpmock.NewStringResponse(http.StatusOK, string(fileContents)), nil
```

Or, if the `httpmock.File()` utility hides errors, switch to explicit Go file IO for fixture files.

