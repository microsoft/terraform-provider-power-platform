# Title

Misuse of `strings.ReplaceAll` for name sanitization

## Path

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_environment_test.go

## Problem

The `strings.ReplaceAll` function is used for sanitizing the `name` field in several test configurations. While this ensures that underscores are replaced in the provided name, it lacks a proper validation mechanism or full sanitization to ensure adherence to naming conventions, particularly if the name includes other unsupported characters.

## Impact

Insufficient sanitization might lead to invalid resource names that could fail API calls or cause unexpected issues in the underlying resource provisioning. This could result in test failures unrelated to the intended functionality being tested. Severity: **Medium**

## Location

Here is the code snippet demonstrating the misuse:

### Code Issue

```go
name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `"
```

## Fix

Replace `strings.ReplaceAll` with a more robust mechanism that validates and sanitizes resource names comprehensively.

### Code Example

```go

import (
    "regexp"
)

// SanitizeName ensures the resource name follows valid naming conventions
func SanitizeName(rawName string) string {
    name := strings.ReplaceAll(rawName, "_", "")
    // Remove all characters that are not alphanumeric, dash, or underscore
    sanitized := regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(name, "")
    return sanitized
}

// Use the new function:
name := SanitizeName(mocks.TestName())

```

This fix replaces unsupported characters and ensures the name strictly adheres to the expected format, reducing the risk of invalid resource names.