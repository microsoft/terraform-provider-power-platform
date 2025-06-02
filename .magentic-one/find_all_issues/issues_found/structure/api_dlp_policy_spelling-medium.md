# Spelling mistake: `covertDlpPolicyToPolicyModelDto` should be `convertDlpPolicyToPolicyModelDto`

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

The function name `covertDlpPolicyToPolicyModelDto` is likely a typo and should be `convertDlpPolicyToPolicyModelDto` to match naming consistency (other conversions are `convert...` as well).

## Impact

Lower readability, confusion for maintainers, decreased discoverability for function. (Severity: Low)

## Location

Used on lines 40, 78, 89, 142

## Code Issue

```go
v, err := covertDlpPolicyToPolicyModelDto(policy)
```

Also other locations:
```go
return covertDlpPolicyToPolicyModel(policy)
```

## Fix

Rename all instances and the function definition itself to `convertDlpPolicyToPolicyModelDto` and `convertDlpPolicyToPolicyModel`.

```go
v, err := convertDlpPolicyToPolicyModelDto(policyDto)
// ... update the actual function names in source and usage accordingly.
```

