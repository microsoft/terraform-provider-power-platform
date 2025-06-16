# Title

Incorrect Spelling in Field Names and Documentation

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages.go

## Problem

Throughout the code and schema definition, the word "Description" is misspelled as "Descprition" and "Applicaiton" instead of "Application". This is present both in struct fields and their documentation/markdown. Having incorrect spelling decreases code readability and causes confusion, especially in APIs, data models, and schema attributes that need to be referenced elsewhere or mapped to upstream/downstream systems.

## Impact

- Low to Medium: 
  - Source code and schema attribute misspellings cause user confusion and difficulty for maintainers.
  - API consumers or integrators might reference incorrect or inconsistent field names, leading to bugs or undocumented behavior.
  - Inconsistent spelling increases chances of runtime errors if fields are referenced dynamically or via reflection.

## Location

- Schema definition under attributes for "application_descprition" and markdown.
- Data model mapping in the Read method and related structs.

## Code Issue

```go
"application_id": schema.StringAttribute{
	MarkdownDescription: "ApplicaitonId",
	Computed:            true,
},
"application_descprition": schema.StringAttribute{
	MarkdownDescription: "Applicaiton Description",
	Computed:            true,
},
...
ApplicationDescprition: types.StringValue(application.ApplicationDescription),
```

## Fix

**Update attribute names, markdown, and struct fields to correct spelling:**

```go
"application_id": schema.StringAttribute{
	MarkdownDescription: "Application ID",
	Computed:            true,
},
"application_description": schema.StringAttribute{
	MarkdownDescription: "Application Description",
	Computed:            true,
},
...
ApplicationDescription: types.StringValue(application.ApplicationDescription),
```

- Update the struct fields/properties to use `ApplicationDescription` consistently.
- Change attribute map keys and MarkdownDescription values to avoid inconsistent spelling.
- Confirm these corrections in all model, mapping, and schema places.
