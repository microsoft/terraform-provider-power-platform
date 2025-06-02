# Title

Verbose Duplicate Imports

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go

## Problem

Imports are generally well-structured, but a review for unnecessary aliases or potentially unused imports (especially for helper/mocks/const packages) should occur as `mocks`, `helpers`, and `constants` are heavily referenced. In large test files, it's also easy to accumulate unused imports as tests evolve.

## Impact

Bloated import section makes finding dependencies slower, and unused imports can trigger linter errors. Severity: low.

## Location

Import block at file top.

## Code Issue

```go
import (
    "fmt"
    "net/http"
    "regexp"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/jarcoal/httpmock"
    "github.com/microsoft/terraform-provider-power-platform/internal/constants"
    "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
    "github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)
```

## Fix

Run `goimports` or `go mod tidy` and manually review the block, removing or consolidating imports. (For instance, some constants or helpers may be unnecessary with config refactor, or mocks could be imported at test package scope.)

```go
// After refactor, block may look like:
import (
    "fmt"
    "net/http"
    "regexp"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/jarcoal/httpmock"
    // Only relevant constant/helpers/mocks kept
)
```
