# Lack of Custom Type for Security Role IDs

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go

## Problem

The `SecurityRoles` field in `UserResourceModel` is defined as a slice of `string` (i.e., `[]string`). If security role IDs have a specific structure or validation requirements, using a plain string misses an opportunity for type safety and reduces self-documentation. Go best practices recommend using custom types for domain-specific identifiers.

## Impact

Using plain strings for identifiers makes code more error-prone, as any string—even unrelated data—may be assigned. This can result in data inconsistency or subtle bugs. Severity: **low** (unless requirements for structure or validation are strict, where it could be **medium**).

## Location

```go
type UserResourceModel struct {
    // ...
    SecurityRoles     []string       `tfsdk:"security_roles"`
}
```

## Code Issue

```go
SecurityRoles     []string       `tfsdk:"security_roles"`
```

## Fix

Define a custom type for SecurityRoleID and use a slice of this type.

```go
type SecurityRoleID string

type UserResourceModel struct {
    // ...
    SecurityRoles     []SecurityRoleID `tfsdk:"security_roles"`
}
```

This allows for centralizing validation and clear intent when handling security role IDs.
