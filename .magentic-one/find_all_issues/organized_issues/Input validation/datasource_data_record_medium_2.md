# Missing Validation on Empty Required Attributes

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

In the schema, fields such as `environment_id` and `entity_collection` are required, but there is no explicit check or validation for these attributes in the `Read` function beyond what Terraform does upstream. If the config retrieves an empty value for these attributes, subsequent calls to the API will fail with possibly less user-friendly errors, or downstream panics if those values are nil.

## Impact

- Weakens guarantees about input data consistency and robustness.
- Downstream errors may be less obvious to the user, resulting in poor UX and harder debugging.
- Severity: **medium**.

## Location

In `Read`, while processing inputs from config/state:

```go
tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", d.FullTypeName()))

resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

if resp.Diagnostics.HasError() {
	return
}

query, headers, err := BuildODataQueryFromModel(&config)
// ... continues
```

## Code Issue

No explicit check for empty required fields:

```go
queryRespnse, err := d.DataRecordClient.GetDataRecordsByODataQuery(ctx, config.EnvironmentId.ValueString(), query, headers)
```

## Fix

Explicitly check that required string attribute values are non-empty when reading and before making API calls:

```go
if config.EnvironmentId.IsUnknown() || config.EnvironmentId.IsNull() || config.EnvironmentId.ValueString() == "" {
	resp.Diagnostics.AddError(
		"Missing Environment ID",
		"The `environment_id` field must be provided and non-empty.",
	)
	return
}
if config.EntityCollection.IsUnknown() || config.EntityCollection.IsNull() || config.EntityCollection.ValueString() == "" {
	resp.Diagnostics.AddError(
		"Missing Entity Collection",
		"The `entity_collection` field must be provided and non-empty.",
	)
	return
}
```
Add similar checks for other required fields as needed. This provides better diagnostics and robustness.
