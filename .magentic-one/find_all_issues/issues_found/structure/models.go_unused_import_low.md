# Title

Unused Import of "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

##

/workspaces/terraform-provider-power-platform/internal/services/solution/models.go

## Problem

The import:

```go
"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
```

could be unnecessary. The actual struct tags use `timeouts.Value`, but if this alias or import isn't properly leveraged or is unnecessary in the context, the import should be cleaned up.

## Impact

Low. Unused imports bloat the code, may lead to confusion, and can trigger linter warnings.

## Location

File-level import block.

## Code Issue

```go
import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	...
)
```

## Fix

Verify the use of `timeouts.Value`—if not actually needed, remove the import; if used, ensure it's clear and documented.

---
