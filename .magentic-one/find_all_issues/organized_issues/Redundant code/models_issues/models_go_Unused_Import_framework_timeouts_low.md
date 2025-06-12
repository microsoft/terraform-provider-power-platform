# Unused Import: github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go

## Problem

The import path `"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"` is referenced, but the actual `timeouts.Value` struct is only used for typing fields, which does not require the entire resource import. In Go, overly broad or redundant imports decrease clarity and may accidentally pull unnecessary code into the build.

## Impact

Increased binary size, slower compile times, and slight readability issue. Severity: **low**.

## Location

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	// ...
)
```

## Code Issue

```go
"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
```

## Fix

If only `timeouts.Value` is required, consider importing only what is needed or, if possible and cleaner, move to a more relevant import. Otherwise, ensure the dependency is necessary and justified.

```go
import (
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
    // ...
)
```

Alternatively (if there is an alias or more granular import available in the package), do:

```go
import timeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/timeouts"
```

and update field typing accordingly.
