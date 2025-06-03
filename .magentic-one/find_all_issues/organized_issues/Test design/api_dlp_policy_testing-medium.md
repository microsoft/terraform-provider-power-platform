# Missing tests for core functions

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

There are no explicit tests or evidence in this file for coverage of (de-)serialization, mapping, or error handling.

## Impact

Potential regressions would not be caught by automated tests; correctness and API integration can't be guaranteed. Severity: Medium

## Location

Whole file

## Code Issue

N/A (testing absence)

## Fix

Add proper unit tests for all core conversion functions and API flows, especially for:

- `convertPolicyModelToDlpPolicy`
- `convertDlpPolicyToPolicyModel`
- Error handling (e.g., 404/not found)
- API integration results

