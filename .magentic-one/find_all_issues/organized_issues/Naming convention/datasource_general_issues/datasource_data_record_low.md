# Naming Consistency Issue for Schema Variable Names

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

Some public schema variable names such as `orderbySchema` and `topSchema` are not consistently named with Go conventions (`orderBySchema`, `topSchema`). The term "order_by" is used in some schema definitions, but the variable uses "orderby" (missing consistency in casing and underscores), which goes against general Go naming conventions and makes code harder to read and search.

## Impact

Minor inconsistencies may cause confusion for maintainers and make it harder to search or refactor code. Severity: **low**.

## Location

Variable names for schema components at the file-level:

```go
var orderbySchema = schema.StringAttribute{
...
}
```

## Code Issue

```go
var orderbySchema = schema.StringAttribute{
	MarkdownDescription: "...",
	//...
}
```

## Fix

Rename the variable to match Go naming conventions and provide consistent variable naming (e.g., `orderBySchema`).

```go
var orderBySchema = schema.StringAttribute{
	MarkdownDescription: "...",
	//...
}
```
Also, update all related usages throughout the file to use the consistent spelling.
