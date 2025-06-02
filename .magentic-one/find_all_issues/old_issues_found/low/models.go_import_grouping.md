### Title

Import Grouping

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

### Problem

The imported libraries are not grouped effectively, mixing third-party and internal imports:

```go
import (
    "context"
    "errors"
    "fmt"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

### Impact

Severity: Low.
Poor grouping affects code readability and compliance with standard Golang practices.

### Location

File: models.go
Top-level import section.

### Code Issue

```go
import (
    "context"
    "errors"
    "fmt"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

### Fix

Group imports into standard sections:

```go
import (
    // Standard Library
    "context"
    "errors"
    "fmt"
    "strings"

    // Third-Party Libraries
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

    // Internal Libraries
    "github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)
```

Enhances readability and standard practices.
