# Title

Hardcoded default values for ProviderVersion, Commit, and Branch

##

/workspaces/terraform-provider-power-platform/common/release.go

## Problem

Currently, the `ProviderVersion`, `Commit`, and `Branch` variables have hardcoded default values set to `"0.0.0-dev"`, `"dev"`, and `"dev"`, respectively. This approach does not provide proper validation or dynamic initialization of these values based on the context or actual release data, potentially leading to issues if these values are not overwritten during the build process.

## Impact

Hardcoded values may cause confusion during debugging or troubleshooting, as they may not correctly represent the actual version, commit, or branch information. If the values are inadvertently not replaced during the build, it could also lead to incorrect metadata being used in the software release. Severity: **medium**

## Location

In the `common` directory, inside the file `release.go`.

## Code Issue

```go
ProviderVersion = "0.0.0-dev" // Default value for development builds
Commit          = "dev"       // Default value for development builds
Branch          = "dev"       // Default value for development builds
```

## Fix

To resolve this issue, you can add validation logic to ensure these values are properly set during runtime or at least provide mechanisms to verify that these values are updated as expected. For example:

```go
package common

import (
	"errors"
	"fmt"
)

var (
	// ProviderVersion is the version of the released provider, set during build/release process with ldflags.
	// This value can not be const as it is set after the build during the linking process.
	ProviderVersion = "0.0.0-dev" // Default value for development builds
	Commit          = "dev"      // Default value for development builds
	Branch          = "dev"      // Default value for development builds
)

// ValidateReleaseDetails checks if release details are properly set
func ValidateReleaseDetails() error {
	if ProviderVersion == "0.0.0-dev" || Commit == "dev" || Branch == "dev" {
		return errors.New("provider release details are not properly initialized; ensure these values are set during build")
	}
	return nil
}

func main() {
	err := ValidateReleaseDetails()
	if err != nil {
		fmt.Println(err)
		// Implement measures to handle invalid release details
	}
}
```

This fix ensures that there is proper validation to catch uninitialized release details during runtime. It alerts developers if the values are not properly overwritten during the build process.
