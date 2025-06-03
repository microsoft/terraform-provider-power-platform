# Description Duplication: MarkdownDescription and Description Are Duplicated

##

/workspaces/terraform-provider-power-platform/internal/modifiers/requires_replace_string_from_non_empty_modifier.go

## Problem

The method `MarkdownDescription` simply calls `Description`, resulting in duplicated descriptions. While this avoids actual duplication, ideally, MarkdownDescription should provide markdown-formatted text, leveraging formatting for better presentation in documentation tools.

## Impact

**Severity: Low**

Not using markdown formatting may result in suboptimal documentation rendering when consumed by markdown-aware tools, which expect rich text formatting to clarify, emphasize, or structure content.

## Location

```go
func (d *requireReplaceStringFromNonEmptyPlanModifier) Description(ctx context.Context) string {
	return "Ensures that change from non empty attribute value will force a replace when changed."
}

func (d *requireReplaceStringFromNonEmptyPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}
```

## Code Issue

```go
func (d *requireReplaceStringFromNonEmptyPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}
```

## Fix

Provide a Markdown-formatted description, leveraging at least basic Markdown for clarity, for example:

```go
func (d *requireReplaceStringFromNonEmptyPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "**Requires replacement**: This forces a resource to be replaced whenever a non-empty attribute value changes."
}
```
This makes the provider documentation clearer in tools that support such formatting.
