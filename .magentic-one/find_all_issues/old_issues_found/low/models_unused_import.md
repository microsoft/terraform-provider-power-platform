# Title

Unused Import `types`

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/models.go

## Problem

The `github.com/hashicorp/terraform-plugin-framework/types` import is declared but never used elsewhere in the code.

## Impact

Unused imports increase the file's load time unnecessarily. It can cause confusion for developers into thinking the imported module is relevant.

Severity: low

## Location

Unused Import Statement

```go
import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

## Code Issue

The `types` import is present but not utilized within the file.

## Fix

Remove the unused import.

```go
import (
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```