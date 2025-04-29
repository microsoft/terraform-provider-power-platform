# Title

Unused imports in the file

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

Multiple imported packages are not used in the program, which unnecessarily adds clutter to the codebase.

## Impact

Unused imports increase the cognitive load when reviewing the code and may confuse developers about utilized features.
Additionally, unnecessary imports can slightly affect compile time. Severity: Medium.

## Location

Line 8 - Import block:
```
"errors"
"fmt"
```

## Code Issue

```go
import (
    "bytes"
    "errors"
    "fmt"
    "net/http"
    "regexp"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/hashicorp/terraform-plugin-testing/knownvalue"
    "github.com/hashicorp/terraform-plugin-testing/statecheck"
    "github.com/hashicorp/terraform-plugin-testing/terraform"
    "github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
    "github.com/jarcoal/httpmock"
    "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
    "github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)
```

## Fix

Remove these unused imports to clean up the code. This helps avoid confusion and makes the file easier to understand and maintain.

```go
import (
    "bytes"
    "net/http"
    "regexp"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/hashicorp/terraform-plugin-testing/knownvalue"
    "github.com/hashicorp/terraform-plugin-testing/statecheck"
    "github.com/hashicorp/terraform-plugin-testing/terraform"
    "github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
    "github.com/jarcoal/httpmock"
    "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
    "github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)
```