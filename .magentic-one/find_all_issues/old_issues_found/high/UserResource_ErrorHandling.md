# Title
Missing Error Handling for `UserClient`

##

`/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go`

## Problem

The structure `UserResource` declares a field `UserClient` but does not include appropriate mechanisms to handle potential errors resulting from operations on this client. This could lead to runtime errors if the client is not properly initialized or fails during usage.

## Impact

Failure to handle errors related to `UserClient` may result in unexpected application crashes or undefined behavior. Severity: High.

## Location

The issue is located within the `UserResource` declaration.

## Code Issue

```go
type UserResource struct {
	helpers.TypeInfo
	UserClient client
}
```

## Fix

Introduce error-handling mechanisms during initialization and usage of `UserClient` to ensure robust operation.

```go
type UserResource struct {
	helpers.TypeInfo
	UserClient client
}

func (r *UserResource) InitializeUserClient(ctx context.Context) error {
    var err error
    r.UserClient, err = initializeClient(ctx)
    if err != nil {
        return fmt.Errorf("failed to initialize UserClient: %v", err)
    }
    return nil
}
```