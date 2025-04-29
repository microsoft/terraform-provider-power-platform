# Title

Lack of meaningful error context in JSON unmarshalling

##

`/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go`

## Problem

When the code encounters an error while unmarshalling the response body, it returns the error directly without adding any context to indicate what failed (e.g., data source, operation, etc.).

## Impact

This results in insufficient debugging information, reducing the observability of errors in the system. The severity is **high** since it complicates debugging and maintenance.

## Location

JSON unmarshalling in `GetLanguagesByLocation`

## Code Issue

```go
err = json.Unmarshal(response.BodyAsBytes, &languages)

if err != nil {
    return languages, err
}
```

## Fix

Wrap unmarshalling errors with additional context for easier debugging.

```go
err = json.Unmarshal(response.BodyAsBytes, &languages)
if err != nil {
    return languages, fmt.Errorf("failed to unmarshal response body: %w", err)
}
```