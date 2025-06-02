# Title

Ambiguity in error handling logic

##

`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go`

---

## Problem

The error handling logic does not consistently differentiate between known, expected errors and unknown, unexpected errors. For example, the differentiation for `customerrors.ERROR_OBJECT_NOT_FOUND` exists but is not applied universally across methods, such as Create, Read, Update, and Delete.

Additionally, custom error codes are being leveraged, but undefined errors seem to be handled without proper contextual messages, reducing error diagnosis's effectiveness.

---

## Impact

The issue is **medium severity** as it affects debugging and could lead to the wrong interpretation of certain error messages or incorrect error handling steps initiated by developers. Although it doesn't compromise code execution directly, it complicates diagnostics and correctness during API and system interaction.

---

## Location

Examples:
- In the `Read` method:
```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
    resp.State.RemoveResource(ctx)
}
```

In other places across the methods (`Configure`, `Update`), errors only leverage generic messages without narrowing scope:
```go
resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
```

---

## Code Issue

### Example Code

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
    resp.State.RemoveResource(ctx)
}
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
```

---

## Fix

Refactor error-handling code to:
1. Consolidate logic for interpreting custom error codes.
2. Update error messages for unexpected errors to include context (e.g., function name, variable state).

### Corrected Code
#### Example in `Read` method:
```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
    resp.State.RemoveResource(ctx)
    return
}
resp.Diagnostics.AddError(
    fmt.Sprintf("Unknown error occurred while reading user resource: %s", err.Error()),
    fmt.Sprintf("Error context: Environment ID: %s, AAD ID: %s", state.EnvironmentId.ValueString(), state.AadId.ValueString()))
)
```

### Explanation

This solution improves debugging by targeted error messages for problematic scenarios, enabling faster problem identification. It also ensures unhandled errors are exposed with sufficient contextual detail.
