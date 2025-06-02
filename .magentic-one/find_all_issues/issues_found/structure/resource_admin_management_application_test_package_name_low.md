# Test Package Name Test Double Exposure

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application_test.go

## Problem

The test file uses a package declaration with a `_test` suffix (`admin_management_application_test`). In Go, this pattern is used to validate the public API of the package from an external perspective, as it forbids access to unexported identifiers of the tested package.

While this approach has correct use-cases (like package-level black-box testing), it is worth noting that if "white-box" testing (access to unexported identifiers, easier refactoring) is desired, not appending `_test` to the package name would be more appropriate.

In this particular context, if tests ever need to access unexported functions or types in `admin_management_application`, they won’t be able to do so unless the package declaration is adjusted.

## Impact

- **Severity: Low**
- May limit flexibility in tests, especially if deeper or more granular testing is identified as necessary later.
- Could potentially force changes in the package’s exported API just for testability.

## Location

```go
package admin_management_application_test
```

## Code Issue

```go
package admin_management_application_test
```

## Fix

Evaluate the need:  
- If only testing the public API as a “black box” is intended, this is fine.
- If you may need “white box” access (unexported symbols), rename the test package to match the implementation package:

```go
package admin_management_application
```
