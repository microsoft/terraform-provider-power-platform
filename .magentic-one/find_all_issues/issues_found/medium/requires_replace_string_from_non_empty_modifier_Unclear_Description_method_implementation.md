# Title

Unclear Description method implementation for `requireReplaceStringFromNonEmptyPlanModifier`

##

`/workspaces/terraform-provider-power-platform/internal/modifiers/requires_replace_string_from_non_empty_modifier.go`

## Problem

The methods `Description` and `MarkdownDescription` provide an identical description string about the purpose of the modifier. However, the implementation lacks extensibility, as future developers may accidentally update one method without realizing the need to update the other, leading to inconsistent method outputs. This redundancy could increase maintenance burden.

## Impact

- **Medium Severity**: This can lead to potential inconsistencies and maintenance challenges, particularly if the descriptions diverge over time.
- Extensibility in the code is reduced due to redundancy.

## Location

```go
func (d *requireReplaceStringFromNonEmptyPlanModifier) Description(ctx context.Context) string {
	return "Ensures that change from non empty attribute value will force a replace when changed."
}

func (d *requireReplaceStringFromNonEmptyPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}
```

## Fix

A more sustainable approach would be to consolidate the logic of the `Description` and `MarkdownDescription` methods, ensuring both are derived from a single source of truth.

```go
func (d *requireReplaceStringFromNonEmptyPlanModifier) Description(ctx context.Context) string {
	description := "Ensures that change from non empty attribute value will force a replace when changed."
	return description
}

func (d *requireReplaceStringFromNonEmptyPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}
```

Alternatively, if the description logic could expand in the future, make it configurable at construction.