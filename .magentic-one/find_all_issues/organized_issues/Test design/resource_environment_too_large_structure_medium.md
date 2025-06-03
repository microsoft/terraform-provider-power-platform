# Title

Code Structure: Large Resource File Mixing Multiple Responsibilities

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go

## Problem

The file contains an extremely large, monolithic implementation for a complex Terraform resource (environment). It contains schema/validation, conversion helpers, CRUD operations, and various helper implementations (such as AI features validation) all within a single file. This makes it difficult to easily test, change, or reason about discrete functionality.

## Impact

- **Severity**: Medium
- Hinders readability and maintainability.
- Increases risk of regression with unrelated changes.
- Discourages targeted automatic or manual testing and code reuse.

## Location

The entire file, but specifically evidence is seen in structure such as:

- `Resource` methods (`Schema`, `Create`, `Update`, `Read`, `Delete`, etc., each large)
- Helper/conversion logic (`convertCreateEnvironmentDtoFromSourceModel`, `updateDataverse`, embedded conversion calls, etc.)
- API and validation logic tightly intermixed with business logic.

## Code Issue

```go
// A long file with methods and implementation for all resource logic and helper/conversion logic
// Example:
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // all update code + helper invocation + conversion
    ...
}
...
// More helpers:
func (r *Resource) updateAllowBingSearch(ctx context.Context, plan *SourceModel) error
...
func addDataverse(ctx context.Context, plan *SourceModel, r *Resource) (string, error)
```

## Fix

Split the file by discrete responsibility. For example:

- Keep the resource struct, CRUD, and config validators in `resource_environment.go`.
- Move all conversion/helper logic (e.g. DTO <-> SourceModel, validation helpers, AI feature validators) into separate files:  
    - `conversion.go`
    - `validators.go`
    - `dataverse_helpers.go`
    - etc.
- Group related helpers/utilities together to promote reuse and testability.

---

**Save as:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_environment_too_large_structure_medium.md`
