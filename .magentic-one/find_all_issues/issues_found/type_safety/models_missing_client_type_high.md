# Use of Unexported Field/Type `client` Without Declaration

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/models.go

## Problem

The struct fields `DataRecordClient client` reference a field or type `client`, but there is no such type declared or imported as `client` in this file. This could be confusing and lead to compilation errors unless it is declared elsewhere in the package with a lower-case name (i.e., as a package-private type).

## Impact

High. If the `client` type is not declared or imported in the package, this code will fail to compile. If it exists elsewhere but is not clearly referenced here, this reduces readability and raises maintainability risks.

## Location

```go
type DataRecordDataSource struct {
	helpers.TypeInfo
	DataRecordClient client
}

// ...

type DataRecordResource struct {
	helpers.TypeInfo
	DataRecordClient client
}
```

## Fix

Make sure to import or declare the `client` type in this file or in the visible package scope. If `client` is a package type, clarify and import as needed:

```go
import "github.com/microsoft/terraform-provider-power-platform/internal/client"

// Then, reference it explicitly:
DataRecordClient client.Client
```

Or, if it is to be declared locally, do so explicitly:

```go
type Client struct { ... }

// use DataRecordClient Client
```
