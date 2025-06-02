# Issue 3

Use of Deprecated/Legacy Httpmock API for File Response

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections_test.go

## Problem

The test code utilizes `httpmock.File("tests/datasource/connections/Validate_Read/get_connections.json").String()` to populate the response for a mocked HTTP call. The idiomatic way to load file contents for use in test responses with `httpmock` is to use `ioutil.ReadFile` or Go's `os`/`io` API directly, and return the file contents.

Using `.String()` on the File handle could produce unexpected behavior, especially across different versions of httpmock or if the method is removed/deprecated.

## Impact

- **Severity:** Low
- May break in the future if the httpmock's API changes or if a minor version update happens
- Test maintainability is reduced and new contributors may not understand this custom approach.

## Location

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/connections/Validate_Read/get_connections.json").String()), nil
```

## Fix

Read the file using standard Go file handling and embed its string contents directly in the mock response.

```go
body, err := ioutil.ReadFile("tests/datasource/connections/Validate_Read/get_connections.json")
if err != nil {
    t.Fatalf("unable to read mock file: %v", err)
}
return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
```
This way, the code is more portable, idiomatic, and will not break if the third-party mocking library changes its semantics.
