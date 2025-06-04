# Test and Mock Issues

This document consolidates issues related to test structure, test data handling, and mock implementations that need to be addressed to improve test maintainability and reliability.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/helpers/hash_test.go`

**Problem:** Use of Magic Strings in File Content

The same string literals like `"same"` and `"different"` are used multiple times as file content. This can lead to errors and inconsistencies if used in multiple places or changed in a single location but not others.

**Impact:** Low severity. This is a minor maintainability issue, but improving this makes content changes less error-prone and increases overall readability.

**Location:**

```go
err := os.WriteFile(file1, []byte("same"), 0644)
...
err = os.WriteFile(file2, []byte("same"), 0644)
...
err = os.WriteFile(file3, []byte("different"), 0644)
```

**Code Issue:**

```go
err := os.WriteFile(file1, []byte("same"), 0644)
```

**Fix:** Declare file content as constants at the beginning of the test:

```go
const (
    sameContent      = "same"
    differentContent = "different"
)
```

Then use:

```go
err := os.WriteFile(file1, []byte(sameContent), 0644)
...
err = os.WriteFile(file3, []byte(differentContent), 0644)
```

---

## Task Completion Instructions

After implementing these fixes:

1. **Run the linter:** `make lint` to ensure code style compliance
2. **Run unit tests:** `make unittest` to verify functionality  
3. **Generate documentation:** `make userdocs` to update auto-generated docs
4. **Add changelog entry:** Use `changie new` to document the changes

**Changie Entry Template:**

```yaml
kind: changed
body: Improved test structure by replacing magic strings with constants for better maintainability
time: [current-timestamp]
custom:
  Issue: "[ISSUE_NUMBER_IF_APPLICABLE]"
```

Replace `[ISSUE_NUMBER_IF_APPLICABLE]` with the relevant GitHub issue number, or remove the custom section if no specific issue exists.
