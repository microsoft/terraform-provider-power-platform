# Title

Missing Test Case for Edge Cases and Error Scenarios

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments_test.go

## Problem

The file only contains tests for successful cases (basic data source retrieval and a mocked API response in unit tests). It does not contain any test covering error cases, such as invalid/malformed environments, invalid API responses (e.g., 500 or 404), empty responses, missing required fields, or API timeouts. Comprehensive test coverage should include scenarios where things go wrong to assure robust provider and resource behavior.

## Impact

Severity: **High**

Without tests covering error and edge cases, regressions or incorrect error handling may go unnoticed. This undermines reliability and makes it harder to guarantee that the provider behaves correctly under adverse or unexpected backend responses (e.g., malformed JSON, network failures).

## Location

- TestAccEnvironmentsDataSource_Basic (no negative tests)
- TestUnitEnvironmentsDataSource_Validate_Read (no negative/edge/error tests)

## Code Issue

```go
func TestAccEnvironmentsDataSource_Basic(t *testing.T) {
    // ... no error/edge tests
}

func TestUnitEnvironmentsDataSource_Validate_Read(t *testing.T) {
    // ... only positive path with HTTP mocks
}
```

## Fix

Add additional test steps and new test functions to cover negative/error and edge scenarios. E.g.:

```go
func TestAccEnvironmentsDataSource_InvalidConfig(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{{
            Config: `data "powerplatform_environments" "all" { invalid_attr = "value" }`,
            ExpectError: regexp.MustCompile(".*unknown attribute.*invalid_attr.*"),
        }},
    })
}

func TestUnitEnvironmentsDataSource_HTTPFailure(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    
    // Simulate a 500 Internal Server Error
    httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/...`, httpmock.NewStringResponder(500, "Internal server error"))
    
    resource.Test(t, resource.TestCase{
        IsUnitTest: true,
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{{
            Config: `data "powerplatform_environments" "all" {}`,
            ExpectError: regexp.MustCompile("500.*Internal server error"),
        }},
    })
}
```

Add similar tests for empty responses, invalid JSON, missing required attributes, etc.
