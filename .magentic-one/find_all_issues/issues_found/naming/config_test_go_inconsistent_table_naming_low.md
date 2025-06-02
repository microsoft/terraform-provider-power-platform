# Inconsistent Naming in Test Table Types

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

The test table struct type for each test is declared as `testData` in all, which is generic and could be more descriptive. Also, naming the struct as `testCase` in the loop (shadowing the table variable) can be confusing.

## Impact

Severity: **Low**  
Minor readability/maintainability issue.

## Location

```go
type testData struct {
    // ...
}
for _, testCase := range []testData{ ... }
```

## Code Issue

```go
type testData struct {
    // ...
}
for _, testCase := range []testData{
    // ...
}
```

## Fix

Rename the struct to a more specific name, e.g. `configStringTestCase`, and use `tc` in the loop:

```go
type configStringTestCase struct {
    // ...
}
for _, tc := range []configStringTestCase{
    // ...
}
```

Repeat for other test tables.
