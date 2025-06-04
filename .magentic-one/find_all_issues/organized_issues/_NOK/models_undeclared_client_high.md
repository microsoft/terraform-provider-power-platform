# Undeclared Type for `client` in DataSource Struct

##
/workspaces/terraform-provider-power-platform/internal/services/capacity/models.go

## Problem

Within the `DataSource` struct, the field `CapacityClient` is declared as `client`, but there is no `client` type imported or defined anywhere in the visible file. This will result in a compilation error, as Go requires all types to be declared or imported. It is unclear from this file whether it should be a pointer, an interface, or a concrete type.

## Impact

High — This is a compilation issue that will prevent the code from building, which is critical for functionality.

## Location

```go
type DataSource struct {
	helpers.TypeInfo
	CapacityClient client
}
```

## Code Issue

```go
	CapacityClient client
```

## Fix

Define or import the correct type for `client`. If it is a struct or interface, ensure it is included in the file or properly imported from another package. For example, if it should be an interface from the current package:

```go
type client interface {
	// method signatures here
}

type DataSource struct {
	helpers.TypeInfo
	CapacityClient client
}
```

Or if it is coming from another package (e.g., `internal/services/capacity/client.go`):

```go
import "github.com/microsoft/terraform-provider-power-platform/internal/services/capacity/client"

type DataSource struct {
	helpers.TypeInfo
	CapacityClient client.Client
}
```

Adjust this to match the actual intended type or import. If it's a pointer (likely), declare as `*client.Client`.
