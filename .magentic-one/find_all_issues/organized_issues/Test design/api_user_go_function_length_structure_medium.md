# Title

Excessive Function Length and Responsibility

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

Several functions (e.g., `RemoveEnvironmentUserSecurityRoles`, `AddEnvironmentUserSecurityRoles`, `CreateDataverseUser`) perform multiple responsibilities and are lengthy, handling API communication, data transformation, logging, and flow control. This leads to poor readability and makes reuse and testing difficult. Best practices recommend smaller, focused functions with clear responsibilities to promote maintainability and understandability.

## Impact

Severity: Medium

Large functions with multiple responsibilities are harder to test, debug, and understand. They increase the risk of hidden bugs and make extending or modifying logic more error-prone.

## Location

For example, `RemoveEnvironmentUserSecurityRoles` (and similar):

## Code Issue

```go
func (client *client) RemoveEnvironmentUserSecurityRoles(ctx context.Context, environmentId, aadObjectId string, securityRoles []string, savedRoles []securityRoleDto) (*userDto, error) {
	...
	userRead, err := client.GetEnvironmentUserByAadObjectId(ctx, environmentId, aadObjectId)
	...
	for _, role := range securityRoles {
		savedRoleData := array.Find(savedRoles, func(roleDto securityRoleDto) bool {
			return roleDto.RoleId == role
		})
		...
	}
	resp, err := client.Api.Execute(...)
	if err != nil {
		return nil, err
	}
	...
	userRead, err = client.GetEnvironmentUserByAadObjectId(ctx, environmentId, aadObjectId)
	...
	user := userDto{
		...
	}
	return &user, nil
}
```

## Fix

Extract independent logic blocks (API formation, executing, processing responses, retry logic) into private helper functions. Example:

```go
func buildRemoveRoleApiUrl(...) string { ... }
func removeRoleFromUser(...) error { ... }
```

Then call these helpers in a concise, high-level workflow within the main function. This reduces cognitive load, enables easier unit testing of sub-parts, and clarifies intent.