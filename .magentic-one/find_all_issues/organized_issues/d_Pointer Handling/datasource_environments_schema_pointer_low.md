# Inconsistent Attribute Pointer Usage in Schema Definition

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments.go

## Problem

In the schema definition section (inside `Schema` method), most attributes are value types. However, for certain attributes such as `billing_policy_id` and `currency_code`, they are created as pointers to `StringAttribute` (e.g., `&schema.StringAttribute{...}`), which is unnecessary and inconsistent since other similar attributes are values and the documentation for Terraform Plugin Framework suggests value type usage unless mutability is required.

## Impact

Unnecessary use of pointers leads to inconsistent code and possible confusion for maintainers; it also may increase risk of accidental nil dereference. **Severity: Low**.

## Location

```go
"billing_policy_id": &schema.StringAttribute{
    MarkdownDescription: "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
    Computed:            true,
},
...
"currency_code": &schema.StringAttribute{
    MarkdownDescription: "Currency name (EUR, USE, GBP etc.)",
    Computed:            true,
},
```

## Code Issue

```go
"billing_policy_id": &schema.StringAttribute{
    MarkdownDescription: "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
    Computed:            true,
},
...
"currency_code": &schema.StringAttribute{
    MarkdownDescription: "Currency name (EUR, USE, GBP etc.)",
    Computed:            true,
},
```

## Fix

Change field assignments to non-pointer value usages, like so:

```go
"billing_policy_id": schema.StringAttribute{
    MarkdownDescription: "Billing policy id (guid) for pay-as-you-go environments using Azure subscription billing",
    Computed:            true,
},
...
"currency_code": schema.StringAttribute{
    MarkdownDescription: "Currency name (EUR, USE, GBP etc.)",
    Computed:            true,
},
```

This makes the code consistent with the rest of the schema and follows the best practices for attribute assignment.
