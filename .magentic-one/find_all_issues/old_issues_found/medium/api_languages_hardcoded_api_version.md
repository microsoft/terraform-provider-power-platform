# Title

Hardcoded API version in query parameters

##

`/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go`

## Problem

The API version (`"2023-06-01"`) is hardcoded in the query parameters. This makes it difficult to update the version across the codebase if it changes and also reduces maintainability.

## Impact

Makes the code harder to maintain and update for future changes. Severity is **medium** because while it won't affect runtime, it could result in issues during updates or changes.

## Location

URL construction in `GetLanguagesByLocation`

## Code Issue

```go
apiUrl.RawQuery = url.Values{
    "api-version": []string{"2023-06-01"},
}.Encode()
```

## Fix

Define the API version as a constant in a centralized location (e.g., within `constants` package or as an API client configuration) and reference it here.

```go
apiUrl.RawQuery = url.Values{
    "api-version": []string{constants.ApiVersion},
}.Encode()
```