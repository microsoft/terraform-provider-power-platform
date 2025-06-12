# Title

Redundant Initialization of Slices in Create and Update Methods

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

In both the `Create` and `Update` resource methods, there are redundant initializations of the `ConnectorGroups` and `CustomConnectorUrlPatternsDefinition` slices within the `dlpPolicyModelDto` struct. After initializing the slice as empty, you immediately overwrite it using the output of append operations. This is unnecessary and may lead to confusion or, if not overwritten, bugs arising from data loss or accidental use of stale slices.

## Impact

Severity: Low

This has a low performance impact and code clarity risk. While not functionally harmful due to the immediate overwrite, it clutters the code and can cause confusion about intentional initialization vs assignment semantics.

## Location

```go
policyToCreate := dlpPolicyModelDto{
	DefaultConnectorsClassification:      plan.DefaultConnectorsClassification.ValueString(),
	DisplayName:                          plan.DisplayName.ValueString(),
	EnvironmentType:                      plan.EnvironmentType.ValueString(),
	Environments:                         []dlpEnvironmentDto{},
	ConnectorGroups:                      []dlpConnectorGroupsModelDto{}, // unnecessary
	CustomConnectorUrlPatternsDefinition: []dlpConnectorUrlPatternsDefinitionDto{}, // unnecessary
}
// ...
policyToCreate.ConnectorGroups = make([]dlpConnectorGroupsModelDto, 0)
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.BusinessGeneralConnectors))
// repeats for \"General\" and \"Blocked\"
```

## Code Issue

```go
policyToCreate := dlpPolicyModelDto{
	DefaultConnectorsClassification:      plan.DefaultConnectorsClassification.ValueString(),
	DisplayName:                          plan.DisplayName.ValueString(),
	EnvironmentType:                      plan.EnvironmentType.ValueString(),
	Environments:                         []dlpEnvironmentDto{},
	ConnectorGroups:                      []dlpConnectorGroupsModelDto{},
	CustomConnectorUrlPatternsDefinition: []dlpConnectorUrlPatternsDefinitionDto{},
}
policyToCreate.Environments = convertToDlpEnvironment(ctx, plan.Environments)
policyToCreate.CustomConnectorUrlPatternsDefinition = convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
policyToCreate.ConnectorGroups = make([]dlpConnectorGroupsModelDto, 0)
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, ...)
```

## Fix

Omit redundant assignments—only assign once and initialize inline or using appends.

```go
policyToCreate := dlpPolicyModelDto{
	DefaultConnectorsClassification: plan.DefaultConnectorsClassification.ValueString(),
	DisplayName:                     plan.DisplayName.ValueString(),
	EnvironmentType:                 plan.EnvironmentType.ValueString(),
}
policyToCreate.Environments = convertToDlpEnvironment(ctx, plan.Environments)
policyToCreate.CustomConnectorUrlPatternsDefinition = convertToDlpCustomConnectorUrlPatternsDefinition(ctx, resp.Diagnostics, plan.CustomConnectorsPatterns)
policyToCreate.ConnectorGroups = []dlpConnectorGroupsModelDto{
	convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.BusinessConnectors), 
	convertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.NonBusinessConnectors),
	convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors),
}
```

---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_dlp_policy_slice_redundancy_low.md
