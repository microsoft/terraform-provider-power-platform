# JSON Tag Case Inconsistency

## /workspaces/terraform-provider-power-platform/internal/services/environment/dto.go

## Problem

In the struct `createTemplateMetadataDto`, the field `PostProvisioningPackages` is mapped to the JSON tag `PostProvisioningPackages`, which uses PascalCase. This is not idiomatic in JSON, where camelCase (or sometimes snake_case) is the convention. Most of the other struct tags use camelCase, so this inconsistency is likely accidental.

## Impact

Such inconsistencies can cause confusion for consumers of the API, especially for those using auto-generated code or documentation. It may also cause issues with tooling expecting standard JSON idioms. Severity: **low**.

## Location

- `createTemplateMetadataDto` struct, around line 241

## Code Issue

```go
PostProvisioningPackages []createPostProvisioningPackagesDto `json:"PostProvisioningPackages,omitempty"`
```

## Fix

Use camelCase for the JSON tag:

```go
PostProvisioningPackages []createPostProvisioningPackagesDto `json:"postProvisioningPackages,omitempty"`
```
