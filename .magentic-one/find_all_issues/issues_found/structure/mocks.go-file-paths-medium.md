# Code Structure: Hardcoded file paths in HTTP mock responders

##

/workspaces/terraform-provider-power-platform/internal/mocks/mocks.go

## Problem

Several HTTP mock responders use hardcoded relative paths for loading JSON files:
```go
httpmock.File("../../services/languages/tests/datasource/Validate_Read/get_languages.json").String()
```
These paths are brittle and could cause failing tests or incorrect behavior if tests are run from a different working directory or the folder structure changes.

## Impact

Hardcoding relative paths reduces test portability, maintainability, and can lead to test breakage when directories change or are run from CI/CD environments. Severity: **medium**

## Location

Multiple HTTP responder registrations, e.g.:
```go
httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\\.bap\\.microsoft\\.com/providers/Microsoft\\.BusinessAppPlatform/locations/(europe|unitedstates)/environmentLanguages\\?api-version=2023-06-01$`),
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/languages/tests/datasource/Validate_Read/get_languages.json").String()), nil
    })
```
and similar.

## Code Issue

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../services/languages/tests/datasource/Validate_Read/get_languages.json").String()), nil
```

## Fix

Use configuration or test environment setup to dynamically resolve file paths, for example by using `os.Getwd()` or setting a test data directory at runtime, or embed the test data files using Go 1.16+ `embed` package for stable, portable pathing.

**Option 1: Use embed for testdata**
```go
import "embed"

//go:embed testdata/get_languages.json
var languagesJSON string

func ActivateEnvironmentHttpMocks() {
    httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`pattern`),
        func(req *http.Request) (*http.Response, error) {
            return httpmock.NewStringResponse(http.StatusOK, languagesJSON), nil
        })
}
```
**Option 2: Compute absolute path using `os.Getwd()`**
```go
func getTestFilePath(relPath string) string {
    wd, _ := os.Getwd()
    return filepath.Join(wd, relPath)
}
...
return httpmock.NewStringResponse(http.StatusOK, httpmock.File(getTestFilePath("../../services/languages/tests/datasource/Validate_Read/get_languages.json")).String()), nil
```
