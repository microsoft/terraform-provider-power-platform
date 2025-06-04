# Inadequate Diagnostic Context for Attribute in Checksum Calculation

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_value_unknown_if_checksum_change_modifier.go

## Problem

When reporting an error in `hasChecksumChanged`, the diagnostic message uses `attribute` (which could be a zero value if unmarshalling fails) rather than the name of the attribute being processed.

## Impact

Error messages do not clearly specify which attribute caused the error, making debugging more difficult. Severity: medium.

## Location

In `hasChecksumChanged`:

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
```

## Fix

Report the name of the attribute, not its value.

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for attribute %q", attributeName), err.Error())
```
