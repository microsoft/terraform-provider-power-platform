# Issue with Markdown Link Syntax in Schema Descriptions

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

The schema attributes use markdown descriptions where the link syntax is incorrect. For example, they use `[...](` instead of standard markdown `[text](url)`. This results in improperly rendered links in documentation and possibly confusion for users reading the descriptions in Terraform and documentation tooling.

## Impact

- Documentation links are broken or not clickable.
- Reduces clarity and professionalism of provider documentation.
- Severity: **low**.

## Location

In lines like:

```go
MarkdownDescription: "Navigation property of the entity collection. \n\nMore information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]",
```

## Code Issue

```go
"More information on (OData Navigation)[https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties]"
```

## Fix

Swap the text and link to proper markdown format:

```go
"More information on [OData Navigation](https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#expand-collection-valued-navigation-properties)"
```
Update all MarkdownDescription fields in the schema to use `[text](url)` format, not `(text)[url]`.
