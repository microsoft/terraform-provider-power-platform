# Unused Import: "github.com/hashicorp/terraform-plugin-framework/types"

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/models.go

## Problem

The import `"github.com/hashicorp/terraform-plugin-framework/types"` appears to only be used for the field `Location types.String` in the `DataSourceModel` struct. If `Location` is not utilized in the provider logic, or if the type should be a base type instead (e.g. `string`), this may be unnecessary. If this usage is intentional and elsewhere, leave as is—otherwise, consider removing to clean dependencies.

## Impact

Unused or unnecessary dependencies increase the binary size, decrease maintainability, and may introduce accidental vulnerabilities or complexity. The severity is **low** unless it is confirmed that the imported type is unneeded.

## Location

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)
```

## Code Issue

```go
	"github.com/hashicorp/terraform-plugin-framework/types"
```

## Fix

If `types.String` is not required (i.e., `Location` could be a bare string, or if the field is not used at all), remove the import and update the field:

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)
```

And in the struct:

```go
	Location string `tfsdk:"location"`
```

If `types.String` is required for Terraform SDK compatibility, no fix is necessary.
