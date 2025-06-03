# No Reference to Unit or Integration Testing

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go

## Problem

The file contains no evidence of in-code testability hooks like interface abstraction or pluggable dependencies for the `AdminManagementApplicationClient`. There are also no references to tests—either as comments, examples, or test scaffolds. This makes it difficult to write unit or integration tests for this resource and may hinder future maintenance or refactoring.

## Impact

Severity: **medium**

Lack of testability is a medium severity concern. If future changes are made, it will become hard to validate correctness or catch regressions without proper test surfaces.

## Location

- Applies to the entire resource file implementation.

## Code Issue

Example (current state):

```go
r.AdminManagementApplicationClient = newAdminManagementApplicationClient(client.Api)
```

There is no way to provide a mock or fake implementation in unit tests.

## Fix

Refactor to inject the client via an interface which can be replaced in test builds:

```go
type ManagementAppClient interface {
    GetAdminApplication(ctx context.Context, id string) (...)
    RegisterAdminApplication(ctx context.Context, id string) (...)
    UnregisterAdminApplication(ctx context.Context, id string) error
}

// In struct
type AdminManagementApplicationResource struct {
    ...
    Client ManagementAppClient
}

// Then inject a mock/fake or the real implementation in production code.
```

Document or scaffold an example test file for the core resource behaviors.
