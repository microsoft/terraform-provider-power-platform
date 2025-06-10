# Title

Potential Inconsistent API Parameter Mapping for Connector Groups

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

The mapping between UI/group labels (“Confidential”, “General”, “Blocked”) and how those map into API-side expectations (`business_connectors`, `non_business_connectors`, etc.) is spread via magic strings throughout the codebase, especially in conversions such as `convertToAttrValueConnectorsGroup` and `convertToDlpConnectorGroup`. Mixing these display group names and schema field references can easily lead to typos and hard-to-update code, as well as mapping mismatches if API requirements or schema names change. The lack of a central mapping function or structure reduces maintainability as the product evolves.

## Impact

Severity: Medium

Medium maintainability and consistency impact. Mistyped or parity-breaking hardcoded values may introduce bugs that are hard to follow, as field semantics are not defined centrally.

## Location

```go
state.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
state.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
state.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)
...
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.BusinessGeneralConnectors))
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.NonBusinessConfidentialConnectors))
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors))
```

## Code Issue

```go
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.BusinessGeneralConnectors))
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.NonBusinessConfidentialConnectors))
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors))
```

## Fix

- Define a central mapping function or constant map for connector group semantics.
- Call central conversion/mapping helpers throughout resource methods.

```go
const (
	BusinessGroup     = \"Confidential\"
	NonBusinessGroup  = \"General\"
	BlockedGroup      = \"Blocked\"
)

var groupFieldMap = map[string]string{
	BusinessGroup:    \"business_connectors\",
	NonBusinessGroup: \"non_business_connectors\",
	BlockedGroup:     \"blocked_connectors\",
}
```

Then use these constants and map lookups in your conversion, resource, and validation logic.

---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_dlp_policy_magic_strings_medium.md
