# Title

Code Duplication in API Response Handling

##

internal/services/solution/api_solution.go

## Problem

In several methods (e.g., `GetSolutionUniqueName`, `GetSolutionById`, `GetSolutions`, `CreateSolution`, `DeleteSolution`, `GetTableData`, `validateSolutionImportResult`), there exists repeated code for handling forbidden and not found HTTP responses right after each `Execute` call:

```go
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```
This repetition across almost every method reduces maintainability and increases the risk of inconsistency if the error handling logic ever changes.

## Impact

Severity: **low**. While this does not present an immediate bug, it decreases maintainability and contributes to code bloat.

## Location

Most functions, e.g.,
```go
resp, err := client.Api.Execute(...)
if err != nil {
    return nil, err
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

## Code Issue

```go
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
```

## Fix

Extract the error handling into a helper function, for example:

```go
func handleCommonApiErrors(api *api.Client, resp *http.Response) error {
    if err := api.HandleForbiddenResponse(resp); err != nil {
        return err
    }
    if err := api.HandleNotFoundResponse(resp); err != nil {
        return err
    }
    return nil
}
```
And then in each method:
```go
if err := handleCommonApiErrors(client.Api, resp); err != nil {
    return nil, err
}
```
