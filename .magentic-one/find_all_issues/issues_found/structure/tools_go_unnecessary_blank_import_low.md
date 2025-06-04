# Unnecessary Blank Import Without Documentation

##

/workspaces/terraform-provider-power-platform/tools.go

## Problem

The code uses a blank import for `github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs` inside the tools package. While blank imports are sometimes needed for tooling dependencies, it is generally best practice to document the reason for this import—especially since it provides context for future maintainers and avoids confusion as to why such import is present.

## Impact

Low severity. The usage of an undocumented blank import can create confusion during maintenance, making it harder to determine if the import is necessary or can be removed, especially for those unfamiliar with the history of the file.

## Location

```go
import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
```

## Code Issue

```go
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
```

## Fix

Add a clear comment explaining why the blank import is used, e.g.:
```go
// Blank import for build tooling: ensures tfplugindocs is included in go.mod for documentation generation.
_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
```

