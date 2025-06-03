# Title

Error messages do not provide enough context for troubleshooting

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

## Problem

Error messages in functions such as `convertFromEnvironmentSettingsModel`, for example `errors.New("failed to convert audit settings to ObjectValue")`, are generic and do not provide sufficient details such as the contents of the object or information about what was expected.

## Impact

Low severity, but it can hinder debugging efforts and reduce the maintainability of the code. Developers might not have enough information to diagnose issues when these errors are reported.

## Location

Examples:

```go
objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
if !ok {
	return nil, errors.New("failed to convert audit settings to ObjectValue")
}
```

## Code Issue

```go
objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
if !ok {
	return nil, errors.New("failed to convert audit settings to ObjectValue")
}
```

## Fix

Include more context in the error message, for example:

```go
objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
if !ok {
	return nil, fmt.Errorf("failed to convert audit settings to ObjectValue, got %T: %+v", auditSettingsObject, auditSettingsObject)
}
```

This provides more information for troubleshooting.
