# Title

Missing Type Validations and Unchecked Nil Returns in Conversion Logic

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

Throughout the resource lifecycle methods (Read, Create, Update), model conversion functions such as `convertToAttrValueConnectorsGroup`, `convertToAttrValueCustomConnectorUrlPatternsDefinition`, and corresponding setters/readers are invoked. There is no type assertion or nil-check logic to handle conversion failures, bad data, or type mismatches. If any conversion returns nil or has a type mismatch, this may cause panics or incorrect assignment to state/plan, and may leave the Terraform state in an inconsistent or corrupted state.

## Impact

Severity: Critical

Unvalidated conversions, unchecked nil errors, or assignment of nil to required schema attributes pose a risk of provider panics and corruption of the user's infrastructure state. This is a critical risk especially during upgrades or with malformed upstream API responses.

## Location

```go
state.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
state.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
state.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)
...
policyToCreate.Environments = convertToDlpEnvironment(ctx, plan.Environments)
policyToCreate.CustomConnectorUrlPatternsDefinition = convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
```

## Code Issue

```go
state.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
...
policyToCreate.CustomConnectorUrlPatternsDefinition = convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
```

## Fix

- Add explicit type checks and nil handling for all conversion functions.  
- If a conversion returns nil or fails, append a diagnostic error rather than assigning nil to required fields.
- Ensure all assignments to required attributes are validated before state updates.

```go
result := convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
if result == nil {
	resp.Diagnostics.AddError("Connector Group Conversion Failed", "Failed to convert 'Confidential' connector group to attribute value.")
	return
}
state.BusinessConnectors = result
```
---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_dlp_policy_missing_nil_checks_critical.md
