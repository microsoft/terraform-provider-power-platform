# Title

Return Pointer to Local Variable (DTO Struct)

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go

## Problem

Returning a pointer to a local variable for a DTO is not strictly unsafe in Go, but it usually warrants explicit attention, especially with larger structs or concurrency.

## Impact

Low. For small structs like DTOs this is commonly acceptable, but it's worth noting as the function API could inadvertently propagate local variable lifetime surprises during future refactoring.

## Location

In each function returning `*adminManagementApplicationDto`.

## Code Issue

```go
var adminApp adminManagementApplicationDto
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &adminApp)

return &adminApp, err
```

## Fix

Ensure this pattern is intentional and add a small comment if kept. For robust code, prefer to clarify by documenting, or for larger structs, perhaps prefer returning by value if suitable.

```go
// Returning pointer to local variable is OK here (small DTO), but document reasoning
return &adminApp, err
```
