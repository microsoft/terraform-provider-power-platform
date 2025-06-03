# Lack of Negative and Edge Test Cases

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings_test.go

## Problem

The current test implementations in the file only verify positive, "happy path" behavior. They do not account for negative cases or edge scenarios such as API failures, malformed responses, or unexpected attributes. Good test coverage should include checks that the provider handles errors and unexpected conditions gracefully and correctly.

## Impact

- **Severity:** Medium
- The lack of negative and edge test scenarios limits the confidence that the provider code and its data source logic behave correctly when things go wrong.
- This can result in silent failures, unhandled panics, or incorrectly handled error states.
- Reduces maintainability and increases chances for regressions when the implementation changes.

## Location

Throughout both test functions in the file.

## Code Issue

```go
func TestAccTenantSettingsDataSource_Validate_Read(t *testing.T) {
    resource.Test(t, resource.TestCase{
        // Only positive/valid config tested here...
    })
}

func TestUnitTestTenantSettingsDataSource_Validate_Read(t *testing.T) {
    // Only valid HTTP mock and correct attributes checked...
}
```

## Fix

Add tests for failure scenarios such as malformed API responses, API errors, missing fields, or unexpected values.

```go
func TestUnitTestTenantSettingsDataSource_Validate_Read_APIFailure(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2023-06-01`,
        func(req *http.Request) (*http.Response, error) {
            // Simulate API error
            return httpmock.NewStringResponse(http.StatusInternalServerError, `{"error": "Internal error"}`), nil
        })

    resource.Test(t, resource.TestCase{
        IsUnitTest: true,
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: `
                data "powerplatform_tenant_settings" "settings" {}`,
                ExpectError: regexp.MustCompile(`Internal error`),
            },
        },
    })
}
```

---

This ensures failures are properly propagated, and error messaging works as expected.
