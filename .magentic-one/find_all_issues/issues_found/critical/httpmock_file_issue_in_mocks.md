# Title

Incorrect Use of `httpmock.File` Method

##

/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go

## Problem

In the `ActivateEnvironmentHttpMocks` function, the `httpmock.File` method is used to fetch mock responses from files. However, it calls the `.String()` method on the `httpmock.File(...)` output, which does not exist in the `httpmock` package. This will lead to a runtime panic when the code is executed.

## Impact

This issue will cause the application to panic during runtime, leading to a disrupted execution. The severity is **critical** as it affects the reliability of the mocks and could block tests or development progress.

## Location

Lines 49 to 61 of the file `/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go`.

## Code Issue

```go
httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/languages/tests/datasource/Validate_Read/get_languages.json").String()), nil

httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/currencies/tests/datasource/Validate_Read/get_currencies.json").String()), nil
```

## Fix

Replace the inappropriate `.String()` method usage with one that works correctly for reading the content of files in `httpmock.File`. For example:

```go
responseData, err := ioutil.ReadFile("../../services/languages/tests/datasource/Validate_Read/get_languages.json")
if err != nil {
    return nil, fmt.Errorf("unable to read mock file: %v", err)
}
httpmock.NewStringResponse(http.StatusOK, string(responseData)), nil

responseData, err = ioutil.ReadFile("../../services/currencies/tests/datasource/Validate_Read/get_currencies.json")
if err != nil {
    return nil, fmt.Errorf("unable to read mock file: %v", err)
}
httpmock.NewStringResponse(http.StatusOK, string(responseData)), nil
```

This fix reads the file content using `ioutil.ReadFile` instead of relying on non-existent functionality and ensures proper error handling.