# Title

Lack of validation for MarkdownDescription in Schema attributes.

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share.go

## Problem

The `MarkdownDescription` field in the Schema attributes (e.g., `timeouts`, `id`, `environment_id`, etc.) is left empty or inconsistent in some sections. Setting proper descriptions helps in user understanding during Terraform usage but the field is unclear in multiple cases.

## Impact

- Reduces the usability of the Terraform provider, causing confusion for users.
- Potentially violates best practices for consistent and descriptive resource documentation.
- **Severity**: _Low impact_, as this issue does not break functionality but reduces clarity.

## Location

```go
resp.Schema = schema.Schema{
	MarkdownDescription: "",
	Attributes: map[string]schema.Attribute{
		"timeouts": ...
		"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the connection share",
				Computed:            true,
		},
		"environment_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the environment",
				Required:            true,
		},
		"connector_name": ...
		...
	},
}
```

## Code Issue

Unclear descriptions in the Schema section:

```go
MarkdownDescription: "",
```

## Fix

Add meaningful descriptions to each MarkdownDescription field. For example:

```go
MarkdownDescription: "Schema definition for the connection_share resource. This includes specifying the required and optional attributes for connection sharing.",
```

For individual attributes:

```go
"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
	MarkdownDescription: "Timeout settings for resource operations such as create, update, delete, and read.",
}),
"id": schema.StringAttribute{
	MarkdownDescription: "Unique identifier for the connection share resource.",
	...
},
"environment_id": schema.StringAttribute{
	MarkdownDescription: "Environment identifier where the connection share is created.",
	 ...
},
```

Provide better user experience through descriptive documentation. Save this change.
