# Title

Non-idiomatic Struct Naming and Field Exporting

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/models.go

## Problem

The struct `AdminManagementApplicationResource` uses an embedded field (`helpers.TypeInfo`) and at least one likely mis-exported field (`AdminManagementApplicationClient client`) with lower case. In Go, exported fields and types should begin with an uppercase letter to adhere to Go conventions and allow access from outside the package. This aids readability, maintainability, and proper package API design.

## Impact

Non-idiomatic naming reduces code readability and portability, limits the struct's use in other packages, and can lead to confusion for developers familiar with Go best practices. Severity: medium

## Location

Line 13, struct: `AdminManagementApplicationResource`

## Code Issue

```go
type AdminManagementApplicationResource struct {
	helpers.TypeInfo
	AdminManagementApplicationClient client
}
```

## Fix

Rename unexported fields/types that need to be accessed outside the package to start with an uppercase letter. Additionally, make sure all embedded fields are exported (public) if they are intended as such.

For example:

```go
type AdminManagementApplicationResource struct {
	helpers.TypeInfo
	Client AdminManagementApplicationClient
}
```

This change allows the field to be accessed externally and adheres to Go's best practices for naming and exporting.
