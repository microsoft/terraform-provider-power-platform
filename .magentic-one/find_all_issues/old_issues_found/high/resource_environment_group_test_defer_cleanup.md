# Issue 4

## Incorrect Cleanup Behavior in HTTPMock

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group_test.go`

## Problem

The `httpmock.DeactivateAndReset()` cleanup function is used, but it's paired with `httpmock.Activate()` without an explicit defer statement, which could lead to test interference if the cleanup is skipped or an error occurs.

## Impact

Such behavior can leave residual mock responses active between tests, leading to flaky tests or unintended test results. Ensuring proper cleanup after each test is critical for test reliability. **Severity: High**

### Location

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Fix

Ensure that all tests are wrapped in functions that correctly use `defer` or dedicated `tearDown` mechanics to enhance cleanup reliability:

```go
func TestUnitEnvironmentGroupResource_Validate_Create(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    
    // Add setup logic here if necessary
    tryFunc := func() {
        resource.Test(t, resource.TestCase{
            IsUnitTest: true,
            ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
            Steps: []resource.TestStep{
                {
                    Config: `
                    resource "powerplatform_environment_group" "test_env_group" {
                        display_name = "test_env_group"
                        description = "test env group"
                    }`,
                    Check: resource.ComposeAggregateTestCheckFunc(
                        resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "display_name", "test_env_group"),
                        resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "description", "test env group"),
                        resource.TestCheckResourceAttrSet("powerplatform_environment_group.test_env_group", "id"),
                        resource.TestCheckResourceAttr("powerplatform_environment_group.test_env_group", "id", "00000000-0000-0000-0000-000000000001"),
                    ),
                },
            },
        })
    }

    tryCatch(tryFunc, tearDown={ /*Add situation-mech statement or guidance*/`
```