# Title

Lack of Error Context in `convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel`

## Path

`/workspaces/terraform-provider-power-platform/internal/services/environment/models.go`

## Problem

The error messages returned in `convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel` lack specific context, which can make debugging difficult. For example, the message `"dataverse object is null or unknown"` does not provide information about the context in which the error occurred or the expected values.

## Impact

- Obscure error messages make debugging and resolving issues more time-consuming.
- Reduces the ability to diagnose issues correctly when running in production.
- Severity: High.

## Location

Defined within the `convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel` function.

## Code Issue

```go
return nil, errors.New("dataverse object is null or unknown")
```

## Fix

Include more context in the error message to aid in debugging.

```go
return nil, fmt.Errorf("dataverse object is null or unknown; context: converting DataverseSourceModel from types.Object in function convertEnvironmentCreateLinkEnvironmentMetadataDtoFromDataverseSourceModel")
```