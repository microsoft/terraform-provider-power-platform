# Function naming: `covertDlpPolicyToPolicyModel` and `covertDlpPolicyToPolicyModelDto`

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

The function `covertDlpPolicyToPolicyModel` is likely a typo of `convertDlpPolicyToPolicyModel`. Consistent and correct naming improves readability and maintainability.

## Impact

Lowers developer experience and codebase quality. (Severity: Low)

## Location

Function and usages throughout file

## Code Issue

```go
func covertDlpPolicyToPolicyModel(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
	...
}
```

## Fix

Rename function definition and all usages to `convertDlpPolicyToPolicyModel`.

```go
func convertDlpPolicyToPolicyModel(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
	...
}
```

