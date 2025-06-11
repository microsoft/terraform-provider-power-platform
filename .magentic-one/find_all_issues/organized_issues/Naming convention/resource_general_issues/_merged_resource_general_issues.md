# Resource General Issues - Merged Issues

## ISSUE 1

# Title

Non-descriptive use of `ctx` and `exitContext` variable names in method bodies

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection.go

## Problem

Within each method, the following pattern is frequently used:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, ...)
defer exitContext()
```

While `ctx` is a standard Go idiom for `context.Context`, the use of `exitContext` as a variable is not self-explanatory - it is actually a cleanup function, not a context type/variable. A more descriptive naming like `cleanup`, `restoreContext`, or `endContextScope` would improve readability.

## Impact

Low: Does not affect correctness, but hurts maintainability/readability for code reviewers unfamiliar with the local idiom. Using more descriptive names helps avoid confusion.

## Location

All methods that use:

## Code Issue

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, ...)
defer exitContext()
```

## Fix

Use a more descriptive variable name for the cleanup function (e.g., `defer cleanup()`):

```go
ctx, cleanup := helpers.EnterRequestContext(ctx, ...)
defer cleanup()
```

This aligns with Go conventions for deferred cleanup and makes intent explicit for future maintainers. Apply for whole codebase


---

## ISSUE 2

# Issue: Function Name Typo - `NewEnterpisePolicyResource`

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go

## Problem

The factory function for this resource is named `NewEnterpisePolicyResource`, which has a typographical error ("Enterpise" should be "Enterprise"). This inconsistency in naming can lead to confusion, reduce code readability, and make it more challenging for other developers to use or search for this constructor function.

## Impact

- Reduces code maintainability and readability.
- Can create difficulty during code search or automated tooling.
- Could cause subtle bugs if the function is mistakenly called elsewhere.
- **Severity:** Low

## Location

Line where the function is defined:

```go
func NewEnterpisePolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "enterprise_policy",
		},
	}
}
```

## Code Issue

```go
func NewEnterpisePolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "enterprise_policy",
		},
	}
}
```

## Fix

Rename the function to `NewEnterprisePolicyResource`, making sure you also update all usage references.

```go
func NewEnterprisePolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "enterprise_policy",
		},
	}
}
```

---

This markdown detailing the naming issue will be saved under:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_enterprise_policy_naming_low.md`


---

## ISSUE 3

# Title

Constant naming does not comply with Go convention

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go

## Problem

The constant `SOLUTION_CHECKER_RULES` uses ALL_CAPS with underscores, which is not idiomatic Go style. Go constants should be named using CamelCase (e.g., `SolutionCheckerRules`). This naming style could cause confusion for contributors who expect Go codebase conventions.

## Impact

Low. This does not affect functionality but may reduce maintainability and code health, especially for new contributors or reviewers familiar with Go best practices.

## Location

At the top of the file:

## Code Issue

```go
const SOLUTION_CHECKER_RULES = "meta-remove-dup-reg, ... web-unsupported-syntax"
```

## Fix

Rename the constant to CamelCase:

```go
const SolutionCheckerRules = "meta-remove-dup-reg, ... web-unsupported-syntax"
```
And update all references accordingly.


---

## ISSUE 4

# Error variable naming inconsistency (`erro`, `errj`) in Create

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

In the `Create` function, you use unconventional error variable names like `erro` and `errj`. For example:

- `originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx)`
- `jsonSettings, errj := json.Marshal(originalSettings)`

This goes against Go naming conventions, which dictate simple, consistent use of `err` for error variables unless there is a strong necessity for disambiguation—typically only in more complex, nested scopes.

## Impact

Reduced code consistency and readability, especially for Go developers who expect the ubiquitous `err`. Severity: low.

## Location

Function `Create`.

## Code Issue

```go
originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx)
if erro != nil {
	resp.Diagnostics.AddError(
		"Error reading tenant settings", fmt.Sprintf("Error reading tenant settings: %s", erro.Error()),
	)
	return
}

jsonSettings, errj := json.Marshal(originalSettings)
if errj != nil {
	resp.Diagnostics.AddError(
		"Error marshalling tenant settings", fmt.Sprintf("Error marshalling tenant settings: %s", errj.Error()),
	)
	return
}
```

## Fix

Use `err` for all error variables. Only use different names in special situations.

```go
originalSettings, err := r.TenantSettingClient.GetTenantSettings(ctx)
if err != nil {
	resp.Diagnostics.AddError(
		"Error reading tenant settings", fmt.Sprintf("Error reading tenant settings: %s", err.Error()),
	)
	return
}

jsonSettings, err := json.Marshal(originalSettings)
if err != nil {
	resp.Diagnostics.AddError(
		"Error marshalling tenant settings", fmt.Sprintf("Error marshalling tenant settings: %s", err.Error()),
	)
	return
}
```

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_tenant_settings.go-error_variable_naming-low.md`


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
