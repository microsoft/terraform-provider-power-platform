# Minor Readability - Redundant Variable Declaration Before Assignment

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

In each test subtest, a variable such as `var configValue basetypes.StringValue` is declared before use. This is fine, but the variable is always set before use to either `NewStringValue` or `NewStringNull`; in Go, short variable declarations (`:=`) could be clearer.

## Impact

Severity: **Low**  
Minor maintainability / readability issue.

## Location

```go
var configValue basetypes.StringValue

if testCase.environmentValue != nil {
    t.Setenv(TEST_ENVIRONMENT_VARIABLE_NAME, *testCase.environmentValue)
}

if testCase.configValue != nil {
    configValue = basetypes.NewStringValue(*testCase.configValue)
} else {
    configValue = basetypes.NewStringNull()
}
```

## Code Issue

```go
var configValue basetypes.StringValue
// then configValue assigned below
```

## Fix

Simplify by merging into one line:

```go
configValue := basetypes.NewStringNull()
if testCase.configValue != nil {
    configValue = basetypes.NewStringValue(*testCase.configValue)
}
```
