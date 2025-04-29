# Title

Improper error handling for `defer response.HttpResponse.Body.Close`

##

`/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go`

## Problem

The code defers the closing of `response.HttpResponse.Body` but does not check if `response` or `response.HttpResponse` is `nil` before attempting to access the `Body`. If `response` is `nil` (e.g., due to a failure in `client.Api.Execute`), the code will panic when `response.HttpResponse.Body.Close()` is executed.

## Impact

This oversight can lead to runtime panics, which makes the application unstable and difficult to debug. Severity is **critical** since it directly affects the stability of the application.

## Location

Function `GetLanguagesByLocation`

## Code Issue

```go
defer response.HttpResponse.Body.Close()
```

## Fix

Check for `nil` before accessing the `Body`. The code should ensure `response` and `response.HttpResponse` are non-nil before attempting to defer closing the `Body`.

```go
if response != nil && response.HttpResponse != nil {
    defer response.HttpResponse.Body.Close()
}
```
