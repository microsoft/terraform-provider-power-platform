# Title

Misspelled function name `covertDlpPolicyToPolicyModelDto`

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem

The function name `covertDlpPolicyToPolicyModelDto` contains a typo (`covert` instead of `convert`). This could be confusing for maintainers and reduces the clarity and discoverability of your code.

## Impact

Low severity. While this doesn’t break code functionality, it can confuse developers and reduce code readability and maintainability due to the inconsistency in naming and possible misspelling.

## Location

Line 16 - Function definition:

```go
func covertDlpPolicyToPolicyModelDto(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
```

## Code Issue

```go
func covertDlpPolicyToPolicyModelDto(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
```

## Fix

Rename the function to correct the spelling, and update all internal references to this function accordingly:

```go
func convertDlpPolicyToPolicyModelDto(policy dlpPolicyDto) (*dlpPolicyModelDto, error) {
```

