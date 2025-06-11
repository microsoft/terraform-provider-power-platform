# Error Handling and Resource Management: Not Checking for `nil` Before Closing Response Body

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The code calls `defer response.HttpResponse.Body.Close()` without ensuring that `response.HttpResponse` or its `Body` field are not nil. If `response.HttpResponse` is nil, a panic may occur.

## Impact

May cause the application to panic at runtime if `HttpResponse` is nil. Severity: **high**.

## Location

```go
defer response.HttpResponse.Body.Close()
```

## Code Issue

```go
defer response.HttpResponse.Body.Close()
```

## Fix

Check for `nil` before deferring the closure:

```go
if response.HttpResponse != nil && response.HttpResponse.Body != nil {
    defer response.HttpResponse.Body.Close()
}
```
