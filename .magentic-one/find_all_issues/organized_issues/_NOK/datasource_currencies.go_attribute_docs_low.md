# Documentation of Data Source Attributes is Not Consistent or Sufficiently Detailed

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go

## Problem

Some documentation strings in schema attributes (e.g., `MarkdownDescription`) are terse or misleading. For example, `"Location of the currencies"` and `"Type of the currency"` do not describe acceptable values, formats, or relationships to the broader Microsoft Power Platform/Currency context.

## Impact

Inadequate documentation reduces usability, increases the support burden, and makes onboarding new contributors or users more difficult. It can increase misconfigurations and lead to support requests or misunderstandings.

**Severity:** Low

## Location

```go
"location": schema.StringAttribute{
	MarkdownDescription: "Location of the currencies",
	Required:            true,
},
...
"type": schema.StringAttribute{
	MarkdownDescription: "Type of the currency",
	Computed:            true,
},
```

## Fix

Expand the `MarkdownDescription` to clarify the meaning, acceptable values, and purpose, e.g.:

```go
"location": schema.StringAttribute{
	MarkdownDescription: "Geographical or logical location identifier for the set of currencies. Should match a valid Dynamics 365 region/location name.",
	Required:            true,
},
...
"type": schema.StringAttribute{
	MarkdownDescription: "The classification type of the currency (such as 'fiat', 'crypto', or other organizational type).",
	Computed:            true,
},
```
