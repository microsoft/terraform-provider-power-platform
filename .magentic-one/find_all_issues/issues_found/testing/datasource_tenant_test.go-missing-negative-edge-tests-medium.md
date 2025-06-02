# Missing Negative and Edge Case Tests

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant_test.go

## Problem

Both `TestUnitTenantDataSource_Validate_Read` and `TestAccTenantDataSource_Validate_Read` only check for the happy path (successful API response and valid attributes). There are no test cases for error conditions such as missing/invalid tenancy info, server errors, malformed responses, or permission errors. Robust testing should include negative and edge scenarios.

## Impact

**Severity: Medium**

- Lack of coverage for possible error states means regressions or improper error handling may slip through unnoticed.
- Confidence in the resilience and reliability of the implementation is limited.

## Location

Both test functions, but especially in:

```go
// Only successful test steps/checks present
Steps: []resource.TestStep{
    {
        Config: `data "powerplatform_tenant" "tenant" {}`,
        Check: resource.ComposeAggregateTestCheckFunc(
            // Attributes...
        ),
    },
},
```

## Fix

Add additional `resource.TestStep` entries to validate:

- API returns 500 or 400 errors
- API returns unexpected or malformed payload
- API returns empty or missing fields
- Required permissions/scopes are missing

Example for a server error (unit test):

```go
// Register a 500 error responder
httpmock.RegisterResponder("GET", ..., 
    func(_ *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusInternalServerError, "server error"), nil
    })

resource.Test(t, resource.TestCase{
    // ...
    Steps: []resource.TestStep{
        {
            Config: `data "powerplatform_tenant" "tenant" {}`,
            ExpectError: regexp.MustCompile("server error"),
        },
        // Existing success case step...
    },
})
```
