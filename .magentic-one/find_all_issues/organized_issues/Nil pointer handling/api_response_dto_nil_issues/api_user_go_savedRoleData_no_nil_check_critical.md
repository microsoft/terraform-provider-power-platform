# Title

Missing Nil Check for `savedRoleData` Result

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

In the `RemoveEnvironmentUserSecurityRoles` function, the variable `savedRoleData` is assigned by calling `array.Find` for each role in `securityRoles`. The code proceeds to use `savedRoleData.Name` without checking whether `savedRoleData` is nil. If `array.Find` does not find a matching role, this will result in a runtime panic (nil pointer dereference).

## Impact

Severity: Critical

If the input lists are inconsistent or data is missing/corrupt, the provider can crash the process, causing cascading failures and overall instability of the Terraform execution and user workflow. This is especially important in infrastructure automation where reliability is paramount.

## Location

Inside RemoveEnvironmentUserSecurityRoles:

## Code Issue

```go
for _, role := range securityRoles {
	savedRoleData := array.Find(savedRoles, func(roleDto securityRoleDto) bool {
		return roleDto.RoleId == role
	})

	remove.Remove = append(remove.Remove, RoleDefinitionDto{
		Id: fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/roleAssignments/%s", environmentId, savedRoleData.Name),
	})
}
```

## Fix

Check for nil before accessing fields and return an appropriate error if `savedRoleData` is not found.

```go
for _, role := range securityRoles {
	savedRoleData := array.Find(savedRoles, func(roleDto securityRoleDto) bool {
		return roleDto.RoleId == role
	})
	if savedRoleData == nil {
		return nil, fmt.Errorf("security role with ID %s not found in savedRoles", role)
	}
	remove.Remove = append(remove.Remove, RoleDefinitionDto{
		Id: fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/roleAssignments/%s", environmentId, savedRoleData.Name),
	})
}
```

This will prevent process panics and provide informative diagnostics to the API caller/user.