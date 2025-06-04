# Title

Error and Warning Handling: Missing Immediate Return After `AddError` or `AddWarning`

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go

## Problem

In several methods (notably, `Read`), after calling `AddError` or `AddWarning` on the diagnostics object, the function does not return immediately. This leads to continued execution, which may act on invalid or incomplete state and possibly cause further, cascading errors.

## Impact

- **Severity**: High
- **Explanation**: Can lead to accessing nil, corrupted, or inconsistent data or duplicate/inconsistent warnings/errors in diagnostics.

## Location

For example, in the `Read` function:

```go
defaultCurrency, err := r.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, envDto.Name)
if err != nil {
    if customerrors.Code(err) != customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND {
        // This is only a warning because you may have BAPI access to the environment but not WebAPI access to dataverse to get currency.
        resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", envDto.Name), err.Error())
    }

    if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
        var dataverseSourceModel DataverseSourceModel
        state.Dataverse.As(ctx, &dataverseSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
        currencyCode = dataverseSourceModel.CurrencyCode.ValueString()
    }
} else {
    currencyCode = defaultCurrency.IsoCurrencyCode
}

var templateMetadata *createTemplateMetadataDto
var templates []string
if !state.Dataverse.IsNull() && !state.Dataverse.IsUnknown() {
    dv, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
    if err != nil {
        resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
        return
    }
    ...
}
```

## Code Issue

```go
dv, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
if err != nil {
    resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
    // should return here, else code acts on a nil dv
}
if dv != nil {
    // ...
}
```

**Several similar cases exist throughout the file.**

## Fix

Add explicit `return` statements immediately after `AddError` or, if appropriate, after critical `AddWarning`s that are expected to not continue processing due to potential invalid state.

```go
dv, err := convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel(ctx, state.Dataverse)
if err != nil {
    resp.Diagnostics.AddError("Error when converting dataverse source model to create link environment metadata", err.Error())
    return
}
if dv != nil {
    // ...
}
```

---

This must be fixed wherever diagnostic error/warning calls are made and continued execution could lead to misuse of the returned data or incorrect resource state.

---
