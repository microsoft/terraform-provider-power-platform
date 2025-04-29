# Title

Inefficient Loop Operation in `GetListStringValues`

##

/workspaces/terraform-provider-power-platform/internal/helpers/config.go

## Problem

The logic within the `GetListStringValues` function iterates twice unnecessarily: first over `environmentVariableNames` to check for environment variable values, and secondly over `defaultValue` if no values are stored from the environment variables. This could lead to performance impacts when dealing with large lists.

## Impact

While not critical, the inefficiency may lead to slower execution in high-performance scenarios or when processing hundreds of list entries. Severity: **Low**

## Location

The function `GetListStringValues`.

## Code Issue

```go
for _, k := range environmentVariableNames {
	if value, ok := os.LookupEnv(k); ok && value != "" {
		values = append(values, types.StringValue(strings.TrimSpace(value)))
	}
}

if len(values) == 0 {
	for _, v := range defaultValue {
		values = append(values, types.StringValue(strings.TrimSpace(v)))
	}
}
```

## Fix

Instead of iterating twice, use a conditional to populate the list in a single loop.

```go
	values := []attr.Value{}

	for i, v := range append(environmentVariableNames, defaultValue...) {
		if i < len(environmentVariableNames) { // Check environment variables
			if value, ok := os.LookupEnv(v); ok && value != "" {
				values = append(values, types.StringValue(strings.TrimSpace(value)))
			}
		} else { // Default values
			values = append(values, types.StringValue(strings.TrimSpace(v)))
		}
	}

	return types.ListValueMust(types.StringType, values)
```