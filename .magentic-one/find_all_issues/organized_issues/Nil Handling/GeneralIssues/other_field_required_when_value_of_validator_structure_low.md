# Unnecessary use of pointer for string type

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

The `otherFieldValue` variable is declared as `new(string)`, resulting in a pointer to a string, but there's no functional need for a pointer since the value is being set directly. This introduces extra indirection that is not idiomatic in Go for basic types unless mutability (across function calls) or explicit nil/no value distinctions are needed. The `GetAttribute` call is also inconsistent in how receiver variables are being used (`currentFieldValue` is a value; `otherFieldValue` is a pointer).

## Impact

Unnecessary use of pointers can make the code harder to read and maintain, and introduces potential for subtle bugs. Severity: **low**.

## Location

```go
otherFieldValue := new(string)
d := req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)
```

## Code Issue

```go
	otherFieldValue := new(string)
	d := req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)
```

## Fix

Just declare as a value variable and pass its address:

```go
	var otherFieldValue string
	d := req.Config.GetAttribute(ctx, paths[0], &otherFieldValue)
```

If you need to check for empty string, just use `otherFieldValue == ""`.

---

This issue impacts code structure and maintainability and should be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/other_field_required_when_value_of_validator_structure_low.md`
