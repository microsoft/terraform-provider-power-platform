# Models Redundant Code Issues

This document consolidates all redundant code issues found in model components of the Terraform Power Platform provider.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go`

### Problem

The import path `"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"` is referenced, but the actual `timeouts.Value` struct is only used for typing fields, which does not require the entire resource import. In Go, overly broad or redundant imports decrease clarity and may accidentally pull unnecessary code into the build.

### Impact

Increased binary size, slower compile times, and slight readability issue. Severity: **low**.

### Location

```go
import (
 "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
 // ...
)
```

### Code Issue

```go
"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
```

### Fix

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

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connection/models.go`

### Problem

The import `github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts` is used in several places for the field type `timeouts.Value`. However, if the surrounding code does not actually use functionality from this import (for example, its fields are simply passed through or never manipulated), the import may be redundant. Similarly, check whether every imported package is necessary.

### Impact

Low. This is a maintainability and tidiness issue. Unused imports do not create runtime errors but may confuse or slow down developers, especially as code evolves.

### Location

- File imports at the top of `/workspaces/terraform-provider-power-platform/internal/services/connection/models.go`

### Code Issue

```go
import (
 "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
 "github.com/hashicorp/terraform-plugin-framework/types"
 "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

### Fix

Remove any unused imports. If all are used, this issue can be ignored. Periodically, running `goimports` or `go fmt` will clean up unnecessary imports:

```go
// Example after removing an unused import
import (
 "github.com/hashicorp/terraform-plugin-framework/types"
 "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

Review each import; if one is unneeded, delete it to tidy up the file.

---

## To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

## Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number

Apply this fix to the whole codebase
