# Title

Lack of Input Validation of the `location` Parameter

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go

## Problem

The `location` parameter is inserted directly into the URL path using `fmt.Sprintf`. There is no checking or sanitization performed. If `location` is empty, malformed, or contains unexpected/slash/control characters, this can lead to malformed URLs and possibly security issues (e.g., path traversal), or functional bugs.

## Impact

**Medium**. Bugs or vulnerabilities can occur if unsanitized or user-generated input is passed.

## Location

Construction of the URL path:

```go
Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location),
```

## Fix

Validate the `location` parameter for empty string and improper characters before using it to build the URL. For example:

```go
if strings.TrimSpace(location) == "" {
    return currencies, fmt.Errorf("location parameter cannot be empty")
}

// Optionally, further sanitize or restrict allowed characters.
```
You might also want to URL-encode the `location` variable (if appropriate).
