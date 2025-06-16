# Inconsistent Use of Required vs Computed Attributes in Schema

##
/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies.go

## Problem

In the schema, some attributes are marked as both `Required` and "`Computed`" (e.g. `"status"` and fields under `"billing_policies"`). Terraform expects either user-provided (`Optional`/`Required`) or provider-computed (`Computed`), but not both: marking a value as both can cause confusion for users and potentially unexpected behavior.

## Impact

"Medium" severity, as it can cause schema validation problems and lead to user confusion or provider malfunction. The fields should adhere to Terraform SDK conventions.

## Location

```go
"status": schema.StringAttribute{
    MarkdownDescription: "The status of the billing policy (Enabled, Disabled)",
    Computed:            true,
    Optional:            true,
},
// ... and similar occurrences in other nested attributes
```

## Code Issue

```go
"status": schema.StringAttribute{
    MarkdownDescription: "The status of the billing policy (Enabled, Disabled)",
    Computed:            true,
    Optional:            true,
},
```

## Fix

Use either `Computed: true` or `Optional/Required: true`, not both. For a data source, values returned from the API should be `Computed: true`, not `Required` or `Optional`.

```go
"status": schema.StringAttribute{
    MarkdownDescription: "The status of the billing policy (Enabled, Disabled)",
    Computed:            true,
},
```

Apply the same change to other attributes that are API-returned and not user-supplied.
