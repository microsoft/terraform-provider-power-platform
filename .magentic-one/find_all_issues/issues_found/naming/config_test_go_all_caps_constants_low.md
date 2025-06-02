# Use of ALL_CAPS Constants for Local Test Data

##

/workspaces/terraform-provider-power-platform/internal/helpers/config_test.go

## Problem

The constants for test environment variable names use ALL_CAPS (e.g., `TEST_ENVIRONMENT_VARIABLE_NAME`) which is idiomatic for exported constants. These are not exported, and Go favors CamelCase for constants unless declaring exported acronyms.

## Impact

Severity: **Low**  
Violates Go naming conventions, minor readbility/naming issue.

## Location

```go
const TEST_ENVIRONMENT_VARIABLE_NAME = "TEST_ENV_VAR"
const TEST_ENVIRONMENT_VARIABLE_NAME1 = "TEST_ENV_VAR_1"
const TEST_ENVIRONMENT_VARIABLE_NAME2 = "TEST_ENV_VAR_2"
```

## Code Issue

```go
const TEST_ENVIRONMENT_VARIABLE_NAME = "TEST_ENV_VAR"
const TEST_ENVIRONMENT_VARIABLE_NAME1 = "TEST_ENV_VAR_1"
const TEST_ENVIRONMENT_VARIABLE_NAME2 = "TEST_ENV_VAR_2"
```

## Fix

Change to Go-style CamelCase and unexported:

```go
const testEnvironmentVariableName = "TEST_ENV_VAR"
const testEnvironmentVariableName1 = "TEST_ENV_VAR_1"
const testEnvironmentVariableName2 = "TEST_ENV_VAR_2"
```
