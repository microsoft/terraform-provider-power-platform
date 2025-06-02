# Title

Incorrect Implementation of Interface: Receiver Should Be Pointer

##

/workspaces/terraform-provider-power-platform/internal/mocks/known_state_value.go

## Problem

The `Check` interface from `github.com/hashicorp/terraform-plugin-testing/knownvalue` is being implemented by `GetKnownValue`. The receiver for `CheckValue` and `String` methods is a value receiver (`GetKnownValue`), but the method mutates the struct (`v.value.Value = otherVal`). This causes the mutation to be made on a copy, not the original struct, effectively dropping the assignment and leading to misuse.

## Impact

High: This leads to silent bugs—mutations appear to succeed but are lost, resulting in incorrect state or test behavior.

## Location

Lines where receiver for methods is `GetKnownValue` instead of `*GetKnownValue`.

## Code Issue

```go
func (v GetKnownValue) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for getKnownValue check, got: %T", other)
	}

	v.value.Value = otherVal

	return nil
}

func (v GetKnownValue) String() string {
	return v.value.Value
}
```

## Fix

Change the method receivers to use a pointer, so that their mutations persist on the target object.

```go
func (v *GetKnownValue) CheckValue(other any) error {
	otherVal, ok := other.(string)
	if !ok {
		return fmt.Errorf("expected string value for getKnownValue check, got: %T", other)
	}
	v.value.Value = otherVal
	return nil
}

func (v *GetKnownValue) String() string {
	return v.value.Value
}
```
