# Title

Unused Imports in datasource_rest_query.go

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

The following import appears to be unused in the provided code:

```go
"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
```

It is imported and used only in calling `timeouts.Attributes`, but if this is the only imported symbol from the whole external package (especially with the unusual import name), code clarity is reduced, and this could be moved to a better, canonical import if needed, or re-evaluated if the package provides only what's needed.

## Impact

Unused or oddly imported packages increase mental load, create confusion in maintenance, or need to be justified with the module's context. While the impact is low, cleaning up imports improves maintainability.

## Location

At the import declaration section.

## Code Issue

```go
"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
```

## Fix

Evaluate whether you need this package, and if it's used only for the `timeouts.Attributes` helper, confirm if your import alias is canonical. If not used at all, remove the import. Otherwise, import as the canonical name or local convention. For example:

```go
// If not used, remove this line
```

Or, if it is used and package is correct, ensure import is as clear as possible.

---

This issue will be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/datasource_rest_query.go_unused_imports_low.md`
