# Lack of Table-Driven Testing and Descriptive Test Names

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant_test.go

## Problem

The tests are written as single blocks, testing only one input/config at a time. Using table-driven tests (a Go best practice) would allow covering multiple scenarios efficiently with clear subtest names and inputs, minimizing code duplication and improving clarity of test cases and report outputs.

## Impact

**Severity: Low**

- Harder to extend to support more cases (e.g., additional configurations).
- Less clarity and maintainability versus table-driven approach.
- Test failure reports are less descriptive (no clear subtest names for various scenarios).

## Location

Throughout the test file, for example:

```go
func TestUnitTenantDataSource_Validate_Read(t *testing.T) {
    // ... only a single test case hard-coded ...
}
```

## Fix

Rewrite tests with a table of scenarios, using `t.Run` for descriptive subtest granularity:

```go
func TestUnitTenantDataSource_Validate_Read(t *testing.T) {
    cases := []struct{
        name    string
        mockSetup func()
        expectErr bool
        // other fields
    }{
        {"ValidResponse", setupValidMock, false},
        {"APIFailure", setupAPIFailureMock, true},
        // add more
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            tc.mockSetup()
            // ... test case logic ...
        })
    }
}
```
