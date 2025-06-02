# Title

Unused Dependency in Code

##

`/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy_test.go`

## Problem

The code imports `fmt` and `github.com/jarcoal/httpmock`, but `fmt` is not used anywhere in the file. This introduces unnecessary code overhead and bloats the file.

## Impact

Unused dependencies increase maintenance complexity and harm code readability. Removing them reduces the chance of introducing accidental bugs tied to these unused imports. Severity: **Low**

## Location

Line 6 (`import "fmt"`) does not serve any purpose in the file since `fmt` is not used.

## Code Issue

```go
import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)
```

## Fix

Remove the unused imports `fmt` and `github.com/jarcoal/httpmock`.

```go
import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)
```
