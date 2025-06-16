# Title

Unused import: github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages.go

## Problem

The import path `github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts` appears to be incorrect or outdated. The correct import path should be `github.com/hashicorp/terraform-plugin-framework/resource/timeouts` for the HashiCorp Terraform plugin SDK. This may be unintentional or a relic from a template, but it may cause build or runtime issues if not corrected.

## Impact

If this code is using a non-existent or deprecated package, it may cause build failures or runtime bugs. If the path is wrong and the package cannot be found, the code will not compile. Severity: **high**.

## Location

```go
import (
    ...
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
    ...
)
...
"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
    Read: true,
}),
```

## Code Issue

```go
"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
```

## Fix

Update the import path to the correct canonical location if using the official SDK, and adjust the call as necessary:

```go
import (
    ...
    "github.com/hashicorp/terraform-plugin-framework/resource/timeouts"
    ...
)
```

And make sure usage of `timeouts.Attributes` and its options stay compatible after the update.
