# Title

Inconsistent JSON Tag Capitalization

##

/workspaces/terraform-provider-power-platform/internal/services/connection/dto.go

## Problem

Some struct field tags use different cases for their JSON keys, for example `NotifyShareTargetOption` (PascalCase) is used in several structs as a JSON tag, while the rest of the tags are in camelCase or all lowercase (e.g., `id`, `type`, `displayName`). This breaks consistency and may lead to errors or confusion when (de)serializing JSON, especially when interacting with external APIs that may expect consistent casing.

## Impact

Medium. Inconsistent JSON key casing can cause subtle bugs if the API expects lowercase/camelCase but the struct produces PascalCase, or vice versa. It may also cause confusion and reduce maintainability if developers do not know which case to use for new fields.

## Location

Several locations, e.g.:

- `permissionPropertiesDto.NotifyShareTargetOption`
- `shareConnectionRequestPutPropertiesDto.NotifyShareTargetOption`
- `shareConnectionResponsePropertiesDto.NotifyShareTargetOption`

## Code Issue

```go
NotifyShareTargetOption string       `json:"NotifyShareTargetOption"`
```

## Fix

Change all occurrences of `NotifyShareTargetOption` JSON tag to use camelCase if that's the convention (as used everywhere else). Update the struct tag accordingly:

```go
NotifyShareTargetOption string       `json:"notifyShareTargetOption"`
```

Apply the same change in every struct that currently uses `NotifyShareTargetOption`.
