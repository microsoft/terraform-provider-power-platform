# Issue: Unused Imports

## File
`/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/models.go`

## Severity
Low

## Problem
The imported package `github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts` is used in the `AdminManagementApplicationResourceModel`, but the `github.com/microsoft/terraform-provider-power-platform/internal/customtypes` and `github.com/microsoft/terraform-provider-power-platform/internal/helpers` appear to be unused in this snippet of code.

## Impact
Unused imports increase the maintenance burden and package size unnecessarily. Developers might assume these imports are required and avoid refactoring or simplifying the code out of caution.

## Location
```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

## Fix
Verify if `customtypes` and `helpers` are intentionally left for future use or can safely be removed:

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	// Removed unused imports if no future use is intended:
	// "github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	// "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```
