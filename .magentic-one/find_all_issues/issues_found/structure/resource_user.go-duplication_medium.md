# Title

Lack of testability and code duplication among CRUD methods

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

The code for resource CRUD operations (Create, Read, Update, Delete) contains duplicated logic for data extraction, diagnostics handling, and context setup. Additionally, the file does not appear to be structured for testability: there is no separation of business logic from resource/plugin framework code, and the tight coupling to the framework/request structure makes isolated unit testing difficult.

## Impact

This reduces maintainability (as updates must be made in multiple places), increases risk of bugs due to missing parallel changes, and discourages automated unit/integration testing. Severity: **Medium**.

## Location

Throughout all main resource methods, e.g.,

- Repeated context/diagnostics setup and handling
- Duplicated plan/state extraction and field assignments
- Similar if/else branches for Dataverse/environment logic

## Code Issue

```go
var plan *UserResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}
// similar blocks in Create/Update/Read, etc.
```

## Fix

Refactor the core business logic for CRUD operations into helper functions or separate service methods that take view-agnostic input/output structures (not tightly coupled to the plugin framework types). Use shared utility methods for extracting plan/state/response and for diagnostics handling. Explore table-driven/unit tests for business logic, making mocking possible by factoring out external calls and isolating side effects.

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_user.go-duplication_medium.md.
