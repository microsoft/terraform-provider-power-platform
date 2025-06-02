# Title

Code Structure: Excessive Function Length and Nesting in Resource Methods

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go

## Problem

Methods such as `Create`, `Read`, and `Update` are extremely lengthy (hundreds of lines) and have deep nesting. They perform multiple types of logic (validation, API preparation, diagnostic error handling, response population, state transformation, etc.) within a single function body. This makes them hard to read, test, and maintain, and increases the risk that changes introduce unintended side effects.

## Impact

- **Severity**: Medium
- Impedes maintainability, comprehension, and ease of onboarding for new contributors.
- Makes function-specific unit testing more difficult or impractical.
- Increases the risk of bugs due to long, complex, and deeply-nested control flow.

## Location

```go
// Example - Create, Read, Update all have this issue:
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // 70+ lines, deep error/nil checks, multiple return points
    ...
}
```

## Code Issue

```go
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // heavy logic, validation, diagnostics, API calling, conversions, more
    ...
}
```

## Fix

Refactor large methods to decompose internal logic into smaller, testable helper functions and minimize maximum indentation level. For example:

- Isolate conversion and API preparation into dedicated functions.
- Move error-diagnostic branching into helpers with semantic naming.
- Group repeated patterns (such as currency/language/result transformation) into concise, reusable helpers.

Example structure:

```go
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    plan := fetchPlan(...)          // Extraction
    envToCreate := prepareEnvDto(...) // Validation & DTO creation
    if err != nil { ... return }
    if err := validateSomething(...); err != nil { ... return }
    // ...
    newState := buildNewState(...)
    resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
```

---

**Save as:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_environment_resource_methods_too_long_medium.md`
