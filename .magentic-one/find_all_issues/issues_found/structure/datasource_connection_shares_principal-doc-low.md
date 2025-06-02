# Documentation: Empty MarkdownDescription for `principal` schema attribute

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

The schema for `principal` leaves `MarkdownDescription` empty.

## Impact

**Severity: low**

Reduces quality and completeness of provider/outbound documentation.

## Location

```go
"principal": schema.SingleNestedAttribute{
	MarkdownDescription: "",
	Computed:            true,
	Attributes: map[string]schema.Attribute{
		// ...
	},
},
```

## Code Issue

```go
"principal": schema.SingleNestedAttribute{
	MarkdownDescription: "",
	Computed:            true,
	Attributes: map[string]schema.Attribute{
		"entra_object_id": schema.StringAttribute{
			MarkdownDescription: "Entra Object Id of the principal",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "Principal Display Name",
			Computed:            true,
		},
	},
},
```

## Fix

Provide a meaningful description:

```go
"principal": schema.SingleNestedAttribute{
	MarkdownDescription: "The principal (user or application) associated with this share.",
	Computed:            true,
	Attributes: map[string]schema.Attribute{
		"entra_object_id": schema.StringAttribute{
			MarkdownDescription: "Entra Object Id of the principal",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "Principal Display Name",
			Computed:            true,
		},
	},
},
```
