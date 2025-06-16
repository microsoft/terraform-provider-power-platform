# Redundant or Misleading MarkdownDescriptions in Schema Attributes

##

/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions.go

## Problem

Within the `Schema` definition, several fields (`modified_time`, `install_time`, `version`) repeat the same MarkdownDescription ("Created time") from the `created_time` attribute, which is semantically incorrect and misleading.

## Impact

* **Medium**. Incorrect documentation can confuse both users and developers, leading to improper use or misunderstanding of the provider's data model. While it does not affect runtime execution, it undermines documentation quality and trust.

## Location

Lines 54-83 within the `Schema` method.

## Code Issue

```go
"created_time": schema.StringAttribute{
	MarkdownDescription: "Created time",
	Computed:            true,
},
"modified_time": schema.StringAttribute{
	MarkdownDescription: "Created time",
	Computed:            true,
},
"install_time": schema.StringAttribute{
	MarkdownDescription: "Created time",
	Computed:            true,
},
"version": schema.StringAttribute{
	MarkdownDescription: "Created time",
	Computed:            true,
},
```

## Fix

Update the `MarkdownDescription` for each attribute to correctly reflect its meaning:

```go
"created_time": schema.StringAttribute{
	MarkdownDescription: "Created time",
	Computed:            true,
},
"modified_time": schema.StringAttribute{
	MarkdownDescription: "Last modified time",
	Computed:            true,
},
"install_time": schema.StringAttribute{
	MarkdownDescription: "Time when the solution was installed",
	Computed:            true,
},
"version": schema.StringAttribute{
	MarkdownDescription: "Solution version",
	Computed:            true,
},
```

