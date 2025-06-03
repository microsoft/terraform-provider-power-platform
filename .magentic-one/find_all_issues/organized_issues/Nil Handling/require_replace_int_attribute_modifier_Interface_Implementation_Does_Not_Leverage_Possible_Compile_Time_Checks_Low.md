# Title

Interface Implementation Does Not Leverage Possible Compile-Time Checks

##

/workspaces/terraform-provider-power-platform/internal/modifiers/require_replace_int_attribute_modifier.go

## Problem

It is idiomatic Go to verify at compile time that a struct implements an interface using a blank identifier assignment:

```go
var _ planmodifier.Int64 = (*requireReplaceIntAttributePlanModifier)(nil)
```

Currently, such an assertion is missing, so if the struct ever stops implementing the required methods for `planmodifier.Int64`, this will only be caught at compile time if explicitly used, instead of immediately failing the build.

## Impact

Low: This is largely a maintainability and future-proofing issue. It helps catch interface implementation regressions early, but does not influence runtime behavior.

## Location

Near the top of the file, after the struct is defined or near the factory function.

## Code Issue

No code line currently exists for this compile-time assertion.

## Fix

Add the following compile-time assertion after the struct or before the factory function:

```go
var _ planmodifier.Int64 = (*requireReplaceIntAttributePlanModifier)(nil)
```

This will ensure that any unimplemented method for the interface will fail at build time.
