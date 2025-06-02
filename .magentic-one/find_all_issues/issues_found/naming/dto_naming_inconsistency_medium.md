# Inconsistent Struct Naming

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/dto.go

## Problem

There is inconsistency in naming patterns for the DTO structs. Some structs end with `Dto` (e.g., `dlpEnvironmentDto`), some with `ModelDto` (e.g., `dlpPolicyModelDto`), and some use `ArrayDto` or similar suffixes. This may confuse maintainers and users about the intended usage or relationship of these types.

For example:

- `dlpPolicyModelDto`, `dlpPolicyDto`, `dlpPolicyDefinitionDto`, `dlpPolicyLastActionDto`, etc.
- `dlpConnectorGroupsModelDto` vs. `dlpConnectorGroupsDto`

## Impact

Severity: **Medium**

Inconsistent naming can reduce maintainability, make refactoring more difficult, and may lead to mistakes in usage, especially as the codebase grows or when tools rely on predictable naming patterns.

## Location

Throughout the file, e.g.:

```go
type dlpPolicyModelDto struct {
    ...
}
type dlpPolicyDto struct {
    ...
}
type dlpConnectorGroupsModelDto struct {
    ...
}
type dlpConnectorGroupsDto struct {
    ...
}
```

## Fix

Adopt a consistent naming convention for DTO struct names. Prefer using a single suffix (e.g., always use `Dto` for Data Transfer Objects). Remove redundant distinctions between `ModelDto`, `Dto`, and similar postfixes unless there is a very clear semantic distinction.

For example:

```go
// Instead of
type dlpPolicyModelDto struct { ... }
type dlpConnectorGroupsModelDto struct { ... }

// Use
type DlpPolicyDto struct { ... }
type DlpConnectorGroupsDto struct { ... }
```

Capitalize struct names if they are exported, and use consistent suffixes for DTOs. Apply the chosen style consistently throughout the file.

