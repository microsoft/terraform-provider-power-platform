# Title

***No Update Function Implemented***

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go`

## Problem

The `Update` function in the resource has no implementation, and instead, it contains the following comment, "Wave features have no updateable attributes". While this might be valid in the current context, leaving methods without implementation but defining them can lead to confusion. This can also be limiting if future updates or additional functionalities are introduced, making it harder to maintain or extend the code.

## Impact

- **Severity:** Medium.
- Risk of missing additions to functionality or confusion among developers working on this project in the future, especially if features change.
- It leads to a lack of maintainability because a better pattern or rationale for handling it (e.g., explicit error or stating it cannot be modified) is not applied.

## Location

**Function Name:**  
`func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)`

## Code Issue

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Wave features have no updateable attributes
}
```

## Fix

The method should explicitly handle the "no updatable attributes" case instead of simply documenting it in comments. For example, returning an error makes it clear to users of the provider that updates are explicitly not allowed.

### Fixed Code Example:

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    tflog.Info(ctx, "Update is not supported for wave features.")
    resp.Diagnostics.AddError(
        "Update Not Supported",
        fmt.Sprintf("%s does not support updates. Consider recreating the resource instead if changes are required.", r.TypeName),
    )
}
```

### Why Fix This Way?
- Makes the behavior of the Update method explicit in runtime (intent is more readable and visible).
- Improves developer and end-user diagnostics.
- Ensures future maintainers can build upon this response pattern, especially if update functionality needs to be introduced.