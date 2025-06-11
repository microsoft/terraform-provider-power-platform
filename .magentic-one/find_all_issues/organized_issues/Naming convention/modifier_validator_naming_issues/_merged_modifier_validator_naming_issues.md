# Modifier Validator Naming Issues - Merged Issues

## ISSUE 1

# Struct type name does not follow Go naming conventions

##
/workspaces/terraform-provider-power-platform/internal/modifiers/force_string_value_unknown_modifier.go

## Problem

The struct `forceStringValueUnknownModifier` does not follow standard Go naming conventions for exported types (should be `ForceStringValueUnknownModifier` if it should be exported). Its unexported status is correct given the current usage (the constructor is exported instead), but the naming could be confusing in larger teams or inconsistent with other code.

## Impact

Deviating from naming conventions can decrease code readability and maintainability. It is a low severity issue but can cause confusion for future maintainers.

## Location

```go
type forceStringValueUnknownModifier struct {
}
```

## Fix

If the type should remain unexported, the naming is acceptable but should be verified for project style consistency. If the type should be exported, rename it:

```go
type ForceStringValueUnknownModifier struct {
}
```

Or clarify with a code comment if unexported is required.


---

## ISSUE 2

# Inconsistent Naming: Use of 'av' as receiver for struct named 'OtherFieldRequiredWhenValueOfValidator'

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

The struct `OtherFieldRequiredWhenValueOfValidator` uses `av` as its receiver shorthand in its methods. This naming is not self-documenting or idiomatic as it does not represent or abbreviate the core words from the struct name. Typical Go style recommends meaningful or at least consistent naming for method receivers (e.g., `v` for validator, `ofrwv` for the initials, or something more descriptive).

## Impact

This makes the code marginally less readable and maintainable, especially for contributors who are unfamiliar with the shorthand. Severity: **low**.

## Location

```go
func (av OtherFieldRequiredWhenValueOfValidator) Description(ctx context.Context) string
func (av OtherFieldRequiredWhenValueOfValidator) MarkdownDescription(_ context.Context) string
func (av OtherFieldRequiredWhenValueOfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse)
func (av OtherFieldRequiredWhenValueOfValidator) Validate(ctx context.Context, req OtherFieldRequiredWhenValueOfValidatorRequest, res *OtherFieldRequiredWhenValueOfValidatorResponse)
```

## Code Issue

```go
func (av OtherFieldRequiredWhenValueOfValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}
```

## Fix

Use a receiver that matches, for example:

```go
func (v OtherFieldRequiredWhenValueOfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}
```

Or abbreviate all initials:

```go
func (ofrwv OtherFieldRequiredWhenValueOfValidator) Description(ctx context.Context) string {
	return ofrwv.MarkdownDescription(ctx)
}
```

---

This issue is a naming issue and should be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/other_field_required_when_value_of_validator_naming_low.md`


---

## ISSUE 3

# Inconsistent or Unclear Variable Naming for Attribute Pairs

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go

## Problem

The variable names `firstAttributePair` and `secondAttributePair` are used to represent slices containing the attribute name and its corresponding checksum attribute name. The naming is vague and can cause confusion about what data is actually stored in these slices. Additionally, storing both pieces as a slice of strings leads to unclear intent and type safety issues; a struct would be more expressive.

## Impact

Unclear naming and lack of structure increase the chance of misusage, make the code less self-documenting, and can confuse maintainers. Severity: Medium.

## Location

```go
func SetBoolValueToUnknownIfChecksumsChangeModifier(firstAttributePair, secondAttributePair []string) planmodifier.Bool {
```

## Code Issue

```go
func SetBoolValueToUnknownIfChecksumsChangeModifier(firstAttributePair, secondAttributePair []string) planmodifier.Bool {
	return &setBoolValueToUnknownIfChecksumsChangeModifier{
		firstAttributePair:  firstAttributePair,
		secondAttributePair: secondAttributePair,
	}
}
...
type setBoolValueToUnknownIfChecksumsChangeModifier struct {
	firstAttributePair  []string
	secondAttributePair []string
}
```

## Fix

Introduce a struct for attribute/checksum pair and use clear naming:

```go
type AttributeChecksumPair struct {
    AttributeName        string
    ChecksumAttributeName string
}

func SetBoolValueToUnknownIfChecksumsChangeModifier(first, second AttributeChecksumPair) planmodifier.Bool {
    return &setBoolValueToUnknownIfChecksumsChangeModifier{
        first:  first,
        second: second,
    }
}

type setBoolValueToUnknownIfChecksumsChangeModifier struct {
    first  AttributeChecksumPair
    second AttributeChecksumPair
}
```


---

## ISSUE 4

# Inconsistent Naming: Struct and Function

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_attribute_unknown_only_if_second_attribute_change.go

## Problem

The struct name `setStringAttributeUnknownOnlyIfSecondAttributeChange` and the function name `SetStringAttributeUnknownOnlyIfSecondAttributeChange` are inconsistent in their capitalization and readability. Go naming convention (as per Effective Go) suggests that exported types and functions should follow CamelCase. Furthermore, the struct's all-lower-case style is non-idiomatic and makes it harder to find types via code search or documentation tooling.

## Impact

- **Maintainability**: Hinders readability for others who may be unfamiliar with the code. CamelCase is the norm for Go struct type names.
- **Discoverability**: Type and function exports become less discoverable via documentation tools.
- **Severity**: Low

## Location

```go
func SetStringAttributeUnknownOnlyIfSecondAttributeChange(secondAttributePath path.Path) planmodifier.String {
	return &setStringAttributeUnknownOnlyIfSecondAttributeChange{
		secondAttributePath: secondAttributePath,
	}
}

type setStringAttributeUnknownOnlyIfSecondAttributeChange struct {
	secondAttributePath path.Path
}
```

## Fix

Rename the struct to follow CamelCase (exported or not). If only used in this file/package, it is fine to keep it unexported but should follow the naming conventions.

```go
type setStringAttributeUnknownOnlyIfSecondAttributeChange struct { // old 
    secondAttributePath path.Path
}

// Suggested improvement
type stringAttributeUnknownIfSecondAttributeChanges struct {
    secondAttributePath path.Path
}
```

And the factory function should reference the new name:

```go
func SetStringAttributeUnknownOnlyIfSecondAttributeChange(secondAttributePath path.Path) planmodifier.String {
	return &stringAttributeUnknownIfSecondAttributeChanges{
		secondAttributePath: secondAttributePath,
	}
}
```


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
