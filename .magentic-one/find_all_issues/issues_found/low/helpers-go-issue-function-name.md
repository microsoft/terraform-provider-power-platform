# Title
Function name is misleading.

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem
The function `covertDlpPolicyToPolicyModelDto` is more correctly identified as `convertDlpPolicyToPolicyModelDto`. The term `covert` may indicate stealth or concealment, causing potential misinterpretation by the developers.

## Impact
Misleading function names can confuse developers, leading to misuse or misunderstanding of the function's purpose. This is considered **low severity** as it does not directly break the code but affects clarity and maintainability.

## Location
Line: 11

## Code Issue

```go
func covertDlpPolicyToPolicyModelDto(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
	policyModel := dlpPolicyModelDto{
```
---

## Fix

```go
// Fix by renaming the function to match its actual purpose:
func convertDlpPolicyToPolicyModelDto(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
	policyModel := dlpPolicyModelDto{
```

Explanation: Renamed function from **covert** to **convert** to better align the name with its behavior.