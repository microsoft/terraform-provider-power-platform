# Title

Unused import statements.

## Path

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/models.go

## Problem

The file contains an unused import statement:

```go
"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
```

Although `timeouts.Value` is referenced in the code, the actual import path `terraform-plugin-framework-timeouts/resource/timeouts` does not exist in the referenced dependency. This suggests an incorrect import path or a mismatch in dependency version.

## Impact

1. **Confusion for developers:** Unused or misconfigured imports imply an oversight that can be confusing for other developers.
2. **Compilation errors:** If the path `terraform-plugin-framework-timeouts/resource/timeouts` is indeed incorrect, the code could result in a compilation error or runtime exception.
3. **Code cleanliness:** Retaining unused imports clutters the file and presents incorrect information about dependencies.

Severity: **High**

## Location

The problem is found at the file's import section:

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)
```

## Code Issue

The issue lies here:

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)
```

## Fix

Confirm the correct path of `timeouts.Value` in the Terraform plugin framework and adjust the import statement accordingly. If the import is entirely unnecessary, remove it.

```go
import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```