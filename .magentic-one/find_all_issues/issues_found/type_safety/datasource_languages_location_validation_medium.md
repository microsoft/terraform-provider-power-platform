# Title

Lack of input validation for `location` attribute in `Read` method

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages.go

## Problem

The `location` attribute from the Terraform state is used directly in API requests without validating if it is empty or malformed before making the call. This can result in unnecessary or erroneous API requests if input is incorrect or not set.

## Impact

This could cause confusing errors or unexpected results from the API, which would be harder to debug by the end user. It is a type safety and input validation issue, with **medium** severity.

## Location

```go
languages, err := d.LanguagesClient.GetLanguagesByLocation(ctx, state.Location.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}
```

## Code Issue

```go
languages, err := d.LanguagesClient.GetLanguagesByLocation(ctx, state.Location.ValueString())
```

## Fix

Add a check to ensure `state.Location.ValueString()` is not empty or invalid before calling the API. Provide a diagnostic error if invalid.

```go
location := state.Location.ValueString()
if location == "" {
    resp.Diagnostics.AddError(
        "Missing or invalid location",
        "The `location` attribute must be set and non-empty.",
    )
    return
}

languages, err := d.LanguagesClient.GetLanguagesByLocation(ctx, location)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}
```
