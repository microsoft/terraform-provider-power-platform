# Empty Update Method With Unclear Purpose

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go

## Problem

The `Update` method on the `Resource` struct is implemented as an empty function with only a comment:

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Wave features have no updateable attributes
}
```
This may be required for compliance with an interface, but it is unclear from the comment or signature. An explicit panic or logging statement could make intent and future maintenance clearer. Alternatively, document clearly that this is an interface compliance stub.

## Impact

An empty method may cause confusion for future maintainers—whether it is a stub, a work in progress, or required for satisfaction of an interface. Severity: **low**.

## Location

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Wave features have no updateable attributes
}
```

## Code Issue

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Wave features have no updateable attributes
}
```

## Fix

Either document clearly that this is required for interface satisfaction, or (preferably) log a message or panic to make it obvious that this code should never be called unless necessary:

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Warn(ctx, "Update called, but wave features have no updateable attributes")
}
```

Or add a comment like:

```go
// Update is present to satisfy the resource.Resource interface; wave features cannot be updated.
```