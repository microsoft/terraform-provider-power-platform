# Title

Unchecked return values and silent error ignoring in strconv.ParseInt

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go

## Problem

In multiple parts of the file, the `strconv.ParseInt` function's error return value is ignored by assigning it to `_`. Ignoring the error may lead to unexpected behaviors if the input string is not a valid integer, leading to unintentional zero values being interpreted and used further in the logic of the provider. This can result in incorrect resource state propagation or masking underlying data/formatting bugs.

## Impact

If the conversion fails, the `maxLimitUserSharing` will be zero and errors will be silently swallowed, potentially producing incorrect Terraform state or misconfigurations without indication to the user or system maintainers. This is a high severity issue because it can corrupt infrastructure state.

## Location

Notably seen in the following code snippet in Create, Update, and Read methods:

## Code Issue

```go
maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
```
And similarly in:
```go
maxLimitUserSharing, _ := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
```

## Fix

Check and handle the error case appropriately and propagate a diagnostic error if parsing fails. For example:

```go
maxLimitUserSharing, err := strconv.ParseInt(env.Properties.GovernanceConfiguration.Settings.ExtendedSettings.MaxLimitUserSharing, 10, 64)
if err != nil {
    resp.Diagnostics.AddError("Error parsing MaxLimitUserSharing as integer", err.Error())
    return
}
```
Add this pattern in Create, Update, and Read methods where parsing is performed.
