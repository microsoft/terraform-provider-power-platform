# Title

Unused Import in `imports` Section

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/models.go`

## Problem

The import `github.com/microsoft/terraform-provider-power-platform/internal/helpers` is not utilized within the file. Unused imports increase the compiled binary size unnecessarily and make the code appear cluttered.

## Impact

**Impact Level**: **Low**

1. **Maintainability**: Extra imports make it harder to quickly identify necessary dependencies.
2. **Performance**: While negligible in small files, unused imports might contribute to larger compiled binary sizes in larger projects.
3. **Readability**: Unused imports introduce noise to the code.

## Location

### Unused Import:
```go
"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
```

## Fix

### Before Removing
```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

### After Removing
Remove the unused import as follows:

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)
```

## Action
Save the above markdown issue in `/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/low/unused_import.md`.