# Title

Inconsistent Error Handling Across Methods

##

/workspaces/terraform-provider-power-platform/internal/services/solution/api_solution.go

## Problem

Error handling across methods is inconsistent. Some methods wrap errors into specific custom errors using `customerrors.WrapIntoProviderError`, while others simply return raw errors. For example, `GetEnvironmentHostById` wraps the error when no environment URL is found, but `GetSolutionById` returns a raw error if no solutions are found.

## Impact

Lack of consistency in error handling can make debugging more difficult and error messages less readable for users. It could result in a product that is harder to maintain and debug.

Severity: low

## Location

The following locations within the file exhibit inconsistent error handling:

- `GetSolutionById`
- `GetSolutionUniqueName`
- `GetEnvironmentHostById`

## Code Issue

```go
if len(solutions.Value) == 0 {
    return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("solution with id '%s' not found", solutionId))
}
...
if environmentUrl == "" {
    return "", customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, please check if the environment has dataverse linked")
}
```

## Fix

Ensure that all methods follow a consistent pattern for error handling, using `customerrors.WrapIntoProviderError` where applicable.

```go
if len(solutions.Value) == 0 {
    return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("solution with id '%s' not found", solutionId))
}
...
if environmentUrl == "" {
    return "", customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, please check if the environment has dataverse linked")
}
```