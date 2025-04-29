# Title

Potential Tight Coupling in Imports

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/models.go`

## Problem

The file imports several packages, including internal helpers (`github.com/microsoft/terraform-provider-power-platform/internal/helpers`) and custom types (`github.com/microsoft/terraform-provider-power-platform/internal/customtypes`) without obvious usage in some sections.

Potential tight coupling may reduce modularity and flexibility for future development. It would be better to limit imports to only those packages actively used within the file.

## Impact

Unnecessary imports can increase dependency complexity, introduce maintainability challenges, and result in larger binary sizes. Severity: **medium**.

## Location

Imports section.

## Code Issue

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

## Fix

Remove unused imports and reorganize the imports section. Example:

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// Removed unused internal packages
)
```