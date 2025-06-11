# Field Name `AadObjectId` and `AadId` Not Consistent

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

There are two fields referencing the Azure Active Directory Object ID:  
- In `userDto`: `AadObjectId string`
- In `UserResourceModel`: `AadId` (used in the conversion function)

This inconsistency in field naming may lead to confusion and bugs, especially in mapping and conversion code. "AAD" is an abbreviation for "Azure Active Directory" and should be used consistently as either `AadObjectId`, `AADObjectId`, or `AadId` everywhere, based on your team/project standards.

## Impact

**Severity: Low**

This impacts readability and maintainability and invites subtle bugs if a new developer mis-maps properties due to the inconsistent naming. It can also make code harder to search and reason about.

## Location

- Field declared as `AadObjectId string` in `userDto`
- Field referenced as `AadId` in `UserResourceModel` and conversion

## Code Issue

```go
type userDto struct {
    ...
    AadObjectId string `json:"azureactivedirectoryobjectid"`
    ...
}

// In the conversion function
AadId:             types.StringValue(userDto.AadObjectId)
```

## Fix

Choose and standardize on a single field name (and abbreviation style) for Azure AD Object ID in all relevant structs. Update codebase-wide accordingly.

```go
type UserDto struct {
    ...
    AadId string `json:"azureactivedirectoryobjectid"`
    ...
}

// Or, always use AadObjectId (if preferred):
type UserDto struct {
    ...
    AadObjectId string `json:"azureactivedirectoryobjectid"`
    ...
}
```

Ensure all code, function, and type names use the chosen standard. If both fields are required for backward compatibility, clearly document why.
