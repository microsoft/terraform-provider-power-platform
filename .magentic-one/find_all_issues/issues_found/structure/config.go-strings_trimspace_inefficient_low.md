# Issue 4: Inefficient Use of Strings.TrimSpace in Loops

##

/workspaces/terraform-provider-power-platform/internal/helpers/config.go

## Problem

Strings are always trimmed with `strings.TrimSpace` for each environment variable value or default string in `GetListStringValues`. If performance or efficiency is a concern, consider if this is necessary for every string, especially if the source is already reliable.

## Impact

This is a low severity issue, with limited scope, as it is a minor inefficiency. However, for environment variables especially, it may occasionally be superfluous unless spaces are expected.

## Location

Line 49-59

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

Document that trimming is intentional (if so), otherwise allow a switch or helper function for performance-sensitive contexts.

```go
// Add a comment:
	// Trimming spaces for environment and default values to ensure clean output.
```

Or, if always expected to be clean, consider omitting `TrimSpace`.
