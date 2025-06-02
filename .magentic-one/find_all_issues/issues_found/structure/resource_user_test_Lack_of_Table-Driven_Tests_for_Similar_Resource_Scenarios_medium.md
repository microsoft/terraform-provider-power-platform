# Title

Lack of Table-Driven Tests for Similar Resource Scenarios

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go

## Problem

Multiple tests with very similar logic and only small configuration changes are defined as nearly copy-paste blocks. This is a missed opportunity for table-driven tests, which would enhance maintainability and condense repetitive test logic.

## Impact

This leads to test duplication and increases the cost of adding, updating, or debugging test scenarios. Severity: medium.

## Location

Blocks where nearly identical test scaffolding is defined, but only config strings or expected values differ. For example, `TestAccUserResource_Validate_Update_Environment_User` and `TestUnitUserResource_Validate_Update_Environment_User` contain repeated config/step blocks differing only in minimal ways.

## Code Issue

```go
func TestAccUserResource_Validate_Update_Environment_User(t *testing.T) {
  resource.Test(t, resource.TestCase{
    ...
    Steps: []resource.TestStep{
      {
        ResourceName: "...",
        Config: "...",
        Check: ...
      },
      {
        ResourceName: "...",
        Config: "...",
        Check: ...
      },
    },
  })
}
```

## Fix

Refactor to use a table of test cases (with configs and expected values) and iterate with a subtest for each, minimizing duplication.

```go
cases := []struct {
    name   string
    config string
    checks resource.TestCheckFunc
}{
    {
        name: "TwoRoles",
        config: "...",
        checks: ..., 
    },
    {
        name: "NoRoles",
        config: "...",
        checks: ..., 
    },
}

for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
        // Use tc.config, tc.checks
    })
}
```
