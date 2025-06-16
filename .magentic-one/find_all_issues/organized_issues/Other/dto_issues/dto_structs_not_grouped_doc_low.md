# Struct Type Definitions Not Grouped or Documented

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

DTO struct types are declared individually and sequentially in a long file, without logical grouping, comments, or section headers. With many DTOs (user, role, environment, principal, etc.) the code becomes harder to scan and to comprehend at a glance. No doc comments or Go style sectioning are present, reducing maintainability and discoverability for consumers.

## Impact

**Severity: Low**

This affects maintainability and team velocity. Onboarding new contributors or debugging issues in DTO definitions becomes harder as the list of types grows. The file also lacks doc comments, making it hard to generate documentation for consumers or consumers of structs in other packages.

## Location

All DTO type declarations. For example:

```go
type userDto struct { ... }
type securityRoleDto struct { ... }
type securityRoleArrayDto struct { ... }
// ...etc...
```

## Code Issue

```go
// Types appear one after another, with no doc comments or sectioning for context.
type userDto struct { ... }
type securityRoleDto struct { ... }
type securityRoleArrayDto struct { ... }
```

## Fix

- Add doc comments to each struct describing its purpose and usage context.
- Use sectioning comments to split the file by business feature ("// User DTOs", "// Role DTOs", etc.).
- Consider grouping related DTOs into sub-files if the DTO module becomes too large.

```go
// UserDto represents a Dataverse user record for serialization.
type UserDto struct {
    // ...
}

// RoleDto represents a security role record.
type RoleDto struct {
    // ...
}

// --- User DTOs ---
// (grouping via comments for quick navigation)
```

This clarifies intent and improves structure as the codebase evolves.
