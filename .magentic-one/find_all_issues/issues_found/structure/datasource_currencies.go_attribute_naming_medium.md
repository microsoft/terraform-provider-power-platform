# Hardcoded and Inconsistent Attribute Names

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go

## Problem

Some attribute names, such as `"code"`, are used in a way that may confuse users or maintainers. For example, `"code"` is described as "Code of the location", but is placed within the `"currencies"` block, not clearly indicating whether this is a currency code or a location code. This kind of ambiguity should be avoided for clarity and future extensibility.

## Impact

Inconsistently named attributes make it harder for users and developers to understand the schema, possibly leading to misconfiguration or misunderstandings about the provider's interface. This is particularly important in public APIs and Terraform providers.

**Severity:** Medium

## Location

```go
"code": schema.StringAttribute{
	MarkdownDescription: "Code of the location",
	Computed:            true,
},
```

## Fix

Rename the field to be explicit, such as `"currency_code"` or `"location_code"`, and update its description accordingly to eliminate ambiguity.

```go
"currency_code": schema.StringAttribute{
	MarkdownDescription: "ISO code of the currency",
	Computed:            true,
},
```

Or

```go
"location_code": schema.StringAttribute{
	MarkdownDescription: "Code of the location associated with the currency",
	Computed:            true,
},
```
