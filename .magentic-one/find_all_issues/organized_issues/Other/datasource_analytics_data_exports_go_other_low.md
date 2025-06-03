# Markdown Issue: Inconsistent Use of Double Line Breaks

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

The Markdown descriptions in the schema definition use multiple `\n\n` line breaks to create spacing between paragraphs, but this is inconsistent and may not render as intended across all Markdown-rendering tools or documentation browsers. Also, some descriptions have a trailing space at the end of text before the newline sequence.

## Impact

- Documentation might be rendered inconsistently
- Minor confusion for consumers reading provider-generated docs
- Aesthetically less polished and slightly harder to maintain

**Severity:** Low

## Location

Example:

```go
MarkdownDescription: "Analytics Data Export configurations. See [documentation](https://learn.microsoft.com/en-us/power-platform/admin/set-up-export-application-insights) for more details.\n\n" +
    "**Note:** This resource is available as **preview**\n\n" +
    "**Known Limitations:** This resource is not supported for with service principal authentication.",
```

## Code Issue

```go
MarkdownDescription: "...\n\n...",
```

## Fix

Review and standardize paragraph breaks according to rendered Markdown requirements in your target documentation generator. For most systems, two newlines is appropriate, but excessive or inconsistent newlines should be avoided. Use one canonical style, e.g.,:

```go
MarkdownDescription: "Analytics Data Export configurations. See [documentation](https://learn.microsoft.com/en-us/power-platform/admin/set-up-export-application-insights) for more details.\n\n**Note:** This resource is available as **preview**\n\n**Known Limitations:** This resource is not supported with service principal authentication.",
```
Remove trailing whitespace before `\n` sequences.
