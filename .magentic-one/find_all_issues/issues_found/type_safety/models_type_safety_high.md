# Title

Unexported Client Type: Use of Undeclared `client` Identifier

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/models.go

## Problem

The structs `DataLossPreventionPolicyResource` and `DataLossPreventionPolicyDataSource` use a field called `DlpPolicyClient client`, but there is no definition or import of a `client` type in this file. If this type is not declared in another imported file or package, this will cause a compilation failure.

## Impact

Severity: High

If the type is not declared or imported, this will prevent the code from compiling and thus block the provider’s build process.

## Location

```go
type DataLossPreventionPolicyResource struct {
	helpers.TypeInfo
	DlpPolicyClient client
}
```

## Code Issue

```go
DlpPolicyClient client
```

## Fix

Make sure to import the correct package and use the correct type name for the client. For example, if the client is in an internal package named `client`, you should import it:
```go
import "github.com/microsoft/terraform-provider-power-platform/internal/client"

// ...

DlpPolicyClient client.Client
```
Or, if the type is already imported but has a different name, ensure consistency.
