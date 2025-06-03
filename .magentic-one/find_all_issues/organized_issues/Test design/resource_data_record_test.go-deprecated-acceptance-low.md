# Use of Deprecated helper/resource Packages in Acceptance Tests

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

The test imports and uses the `github.com/hashicorp/terraform-plugin-testing/helper/resource` package, and related sub-packages (`terraform`, `statecheck`, etc.). The Terraform Plugin Testing Framework documentation strongly encourages moving to the newer `testframework` APIs and `terraform-plugin-testing`'s V6 approach, as the old `helper/resource` APIs are gradually being deprecated.

## Impact

- **Obsolescence**: Risk of breaking tests or lacking features/bugfixes as upstream support wanes.
- **Upgrade Friction**: Harder transition in the future if API is eventually removed.
- **New Contributor Confusion**: Contributors may expect or look for the V6 APIs.
- **Severity**: Low

## Location

Imports list:

```go
"github.com/hashicorp/terraform-plugin-testing/helper/resource"
"github.com/hashicorp/terraform-plugin-testing/terraform"
"github.com/hashicorp/terraform-plugin-testing/statecheck"
// etc.
```

and throughout test logic, e.g.:

```go
resource.Test(t, resource.TestCase{
	// ...
})
```

## Fix

Refactor to use only the documented and supported `terraform-plugin-testing/v6` APIs. Prefer `github.com/hashicorp/terraform-plugin-testing/v6/helper/resource` and related modules, and transition away from `statecheck`/`terraform` in new and existing tests.

Consult the latest migration guides from HashiCorp for examples and migration steps.

---
