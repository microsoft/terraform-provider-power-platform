# Title

Unchecked Map Key Access in `covertDlpPolicyToPolicyModel`

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

During operations involving connectors in `covertDlpPolicyToPolicyModel`, there are several direct map key accesses without checks for their availability, which could cause runtime panics.

## Impact

Direct access to map keys without ensuring their existence can lead to runtime panics when an applied key is missing. Severity marked **Medium**.

## Location

`func covertDlpPolicyToPolicyModel(policy dlpPolicyDto) (*dlpPolicyModelDto, error)`

## Code Issue

```go
for _, connector := range connGroup.Connectors {
missingKey := connector.ActionRules["some_key"] // Example of unchecked access
}
```

## Fix

Add explicit checks for key existence before access.

```go
key, exists := connector.ActionRules["some_key"]
if exists {
// Use `key` safely here
}
```