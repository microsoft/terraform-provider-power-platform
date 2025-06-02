# Title

Unused Imports

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/models.go

## Problem

The import `timeouts` from the `terraform-plugin-framework-timeouts/resource` is not referenced in the rest of the file.

## Impact

Unused imports clutter code and may mislead developers into thinking they are required. Severity is "low" as unused imports do not lead to runtime issues but impact code cleanliness.

## Location

Import statement at line 4.

## Code Issue

```go
import (
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

## Fix

Remove the `timeouts` import if it is unused and does not serve any current purpose.

```go
import (
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```