# Title

Repetitive Calls with Potential Performance Issue

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages.go

## Problem

The Read function iterates over `applications` and, for each iteration, calls `state.Name.ValueString()` and `state.PublisherName.ValueString()` for each comparison, leading to potentially redundant work if those methods are not simple value fetches (e.g., if there's underlying validation/conversion).

## Impact

- Low (in this specific context, performance impact is small due to small dataset).
- If the ValueString() accessor has any logic or computation, redundancy could have a small impact as the dataset grows.
- Cleaner code and slight optimization available by extracting these values before the loop.

## Location

- Read method, inside for loop.

## Code Issue

```go
for _, application := range applications {
	if (state.Name.ValueString() != "" && state.Name.ValueString() != application.ApplicationName) ||
		(state.PublisherName.ValueString() != "" && state.PublisherName.ValueString() != application.PublisherName) {
		continue
	}
	// ...
}
```

## Fix

**Fetch the values once before the loop:**

```go
name := state.Name.ValueString()
publisherName := state.PublisherName.ValueString()

for _, application := range applications {
	if (name != "" && name != application.ApplicationName) ||
		(publisherName != "" && publisherName != application.PublisherName) {
		continue
	}
	// ...
}
```
