# Title

Missing State Reinitialization Before Accumulating Applications List

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages.go

## Problem

In the Read method, `state.Applications` is potentially appended to on each Read. Since Terraform framework may preserve state between calls (depending on usage and errors), not clearing or reinitializing this list can result in duplicated or stale data if Read is called multiple times without resetting state.

## Impact

- Medium: Could cause duplicated entries in the returned list upon multiple reads.
- Data consistency issues for users.

## Location

- Read method, before iterating and appending to `state.Applications`.

## Code Issue

```go
for _, application := range applications {
	// ...
	state.Applications = append(state.Applications, app)
}
```

## Fix

**Explicitly reinitialize `state.Applications` before appending new values:**

```go
state.Applications = []TenantApplicationPackageDataSourceModel{}

for _, application := range applications {
	// ...
	state.Applications = append(state.Applications, app)
}
```
