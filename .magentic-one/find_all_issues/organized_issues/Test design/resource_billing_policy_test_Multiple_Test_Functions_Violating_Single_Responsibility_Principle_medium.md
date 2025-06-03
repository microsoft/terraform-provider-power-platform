# Multiple Test Functions Violating Single Responsibility Principle

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

The test file includes multiple large test functions, each testing both positive and negative cases, API mock registration, and resource checks in a single function. This makes the tests harder to follow, debug, and maintain. They violate the Single Responsibility Principle for unit tests by not separating test scenarios.

## Impact

- **Maintainability**: Hard to update or reason about failures.
- **Readability**: Increases cognitive load for contributors and reviewers.
- **Extensibility**: Adding or modifying one scenario risks breaking unrelated scenarios.

**Severity: Medium**

## Location

Throughout the file, especially in functions such as:

```go
func TestUnitTestBillingPolicyResource_Validate_Create(t *testing.T) { ... }
func TestUnitTestBillingPolicy_Validate_Update(t *testing.T) { ... }
func TestUnitTestBillingPolicy_Validate_Update_ForceRecreate(t *testing.T) { ... }
func TestUnitTestBillingPolicy_Validate_Create_WithoutFinalStatusInPostResponse(t *testing.T) { ... }
```

## Code Issue

```go
func TestUnitTestBillingPolicyResource_Validate_Create(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    mocks.ActivateEnvironmentHttpMocks()

    httpmock.RegisterResponder(...) // several responders

    resource.Test(t, resource.TestCase{
        ...
        Steps: []resource.TestStep{
            {
                Config: `...`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // multiple attr checks
                ),
            }
        },
    })
}
```

## Fix

Split each broad test into smaller, focused tests dedicated to a single scenario or outcome. Register only the relevant mocks within those smaller tests.

```go
func TestBillingPolicyResource_Create_Success(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Setup only the POST and GET responders for the successful creation case
    httpmock.RegisterResponder("POST", "...", ...)
    httpmock.RegisterResponder("GET", "...", ...)

    resource.Test(t, resource.TestCase{
        IsUnitTest: true,
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{ ... },
    })
}

func TestBillingPolicyResource_Create_ErrorResponse(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Setup POST responder that returns an error for this scenario
    httpmock.RegisterResponder("POST", "...", ...)

    resource.Test(t, resource.TestCase{
        IsUnitTest: true,
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{ ... },
    })
}
```
