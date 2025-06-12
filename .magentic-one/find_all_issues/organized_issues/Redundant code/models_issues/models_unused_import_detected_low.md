# Unused Import Detected

##

/workspaces/terraform-provider-power-platform/internal/services/connection/models.go

## Problem

The import `github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts` is used in several places for the field type `timeouts.Value`. However, if the surrounding code does not actually use functionality from this import (for example, its fields are simply passed through or never manipulated), the import may be redundant. Similarly, check whether every imported package is necessary.

## Impact

Low. This is a maintainability and tidiness issue. Unused imports do not create runtime errors but may confuse or slow down developers, especially as code evolves.

## Location

- File imports at the top of `/workspaces/terraform-provider-power-platform/internal/services/connection/models.go`

## Code Issue

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

## Fix

Remove any unused imports. If all are used, this issue can be ignored. Periodically, running `goimports` or `go fmt` will clean up unnecessary imports:

```go
// Example after removing an unused import
import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

Review each import; if one is unneeded, delete it to tidy up the file.
