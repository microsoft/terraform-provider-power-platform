# Title

Direct state mutation on `Value` field without preallocation or guarding

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages.go

## Problem

In the `Read` method, the code appends to the `state.Value` slice without first resetting it or ensuring it is initialized. If the `Read` operation is called multiple times during the provider lifecycle (e.g., a refresh), the slice could be appended to repeatedly, causing duplicate or stale entries in state.

## Impact

This can result in inconsistent or duplicate data returned to Terraform, potential memory bloat or state drift. This is a **medium severity** resource management issue.

## Location

```go
for _, language := range languages.Value {
	state.Value = append(state.Value, DataModel{
		ID:              language.ID,
		Name:            language.Name,
		DisplayName:     language.Properties.DisplayName,
		LocalizedName:   language.Properties.LocalizedName,
		LocaleID:        language.Properties.LocaleID,
		IsTenantDefault: language.Properties.IsTenantDefault,
	})
}
```

## Code Issue

```go
for _, language := range languages.Value {
	state.Value = append(state.Value, DataModel{
		ID:              language.ID,
		Name:            language.Name,
		DisplayName:     language.Properties.DisplayName,
		LocalizedName:   language.Properties.LocalizedName,
		LocaleID:        language.Properties.LocaleID,
		IsTenantDefault: language.Properties.IsTenantDefault,
	})
}
```

## Fix

Reset the `state.Value` slice before appending:

```go
state.Value = make([]DataModel, 0, len(languages.Value))
for _, language := range languages.Value {
	state.Value = append(state.Value, DataModel{
		ID:              language.ID,
		Name:            language.Name,
		DisplayName:     language.Properties.DisplayName,
		LocalizedName:   language.Properties.LocalizedName,
		LocaleID:        language.Properties.LocaleID,
		IsTenantDefault: language.Properties.IsTenantDefault,
	})
}
```
This ensures the state always reflects the latest retrieved data without duplication.
