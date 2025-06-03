# Duplication of Responder Registrations Reduces Maintainability (Code Structure/Maintainability)

##

/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install_test.go

## Problem

There are multiple `RegisterResponder` calls with the same URL pattern and logic, both in the same and multiple test functions. This leads to significant code duplication across the unit tests, increasing maintenance effort and the risk of inconsistencies if one copy is changed and others are not.

## Impact

Severity: Medium

This increases cognitive load for maintainers, duplication will lead to bugs if only one copy is updated, and the test file becomes harder to read and maintain. A single function or map declaring common responders would be easier to maintain and keep consistent.

## Location

Duplicate responders in:

```go
httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`, ...)
```

...and similar patterns in several test functions.

## Code Issue

```go
httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
    func(req *http.Request) (*http.Response, error) {
        id := httpmock.MustGetSubmatch(req, 1)
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Install/get_environment_%s.json", id)).String()), nil
    },
)
```

## Fix

Define a function or helper to encapsulate the registration for repeated responder logic, and call it from each test as needed.

```go
// Helper function for registering environment GET mocks
func registerEnvironmentGetResponder(namespace string) {
    httpmock.RegisterResponder("GET", getEnvPattern, func(req *http.Request) (*http.Response, error) {
        id := httpmock.MustGetSubmatch(req, 1)
        file := fmt.Sprintf("services/environment/tests/resource/%s/get_environment_%s.json", namespace, id)
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File(file).String()), nil
    })
}

// Usage:
registerEnvironmentGetResponder("Validate_Install")
registerEnvironmentGetResponder("Validate_No_Dataverse")
```
