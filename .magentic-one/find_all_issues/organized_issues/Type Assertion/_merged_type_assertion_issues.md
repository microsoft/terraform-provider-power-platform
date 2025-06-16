# Type Assertion Issues in Terraform Provider Power Platform

This document consolidates all type assertion related issues found in the codebase that need to be addressed to improve type safety and prevent runtime failures.

## ISSUE 1

### Equality Method Type Assertion Can Fail for Pointer Receiver

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_value.go

### Problem

In the `Equal` method, the code attempts to cast the provided `attr.Value` to `UUIDValue` with the following assertion:

```go
other, ok := o.(UUIDValue)
```

However, this will fail if `o` is a pointer to a `UUIDValue`. In Go, if the value passed is of type `*UUIDValue`, this assertion will return false, potentially causing equality checks to fail unexpectedly. Since Terraform frameworks often pass around values as either structs or pointers, using a pointer receiver (or handling both cases) is recommended.

### Impact

**Medium Severity**: Incorrect equality checks may result in subtle bugs where UUIDs that are ostensibly the same are not considered equal, causing resource drift detection and other logic relying on equality to malfunction.

### Location

Method `Equal` in `UUIDValue`

### Code Issue

```go
func (v UUIDValue) Equal(o attr.Value) bool {
    other, ok := o.(UUIDValue)
    if !ok {
        return false
    }

    return v.StringValue.Equal(other.StringValue)
}
```

### Fix

Change to support both `UUIDValue` and `*UUIDValue` for comparison, or always use pointers consistently throughout your codebase. Here's a recommended fix:

```go
func (v UUIDValue) Equal(o attr.Value) bool {
    switch other := o.(type) {
    case UUIDValue:
        return v.StringValue.Equal(other.StringValue)
    case *UUIDValue:
        if other == nil {
            return false
        }
        return v.StringValue.Equal(other.StringValue)
    default:
        return false
    }
}
```

---

Apply this fix to the whole codebase

## To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

## Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
