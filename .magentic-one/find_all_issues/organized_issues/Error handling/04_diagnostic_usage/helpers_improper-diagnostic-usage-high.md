# Title

Improper Diagnostic Usage for Error Handling in `convertToDlpConnectorGroup` and `convertToDlpCustomConnectorUrlPatternsDefinition`

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem

In both `convertToDlpConnectorGroup` and `convertToDlpCustomConnectorUrlPatternsDefinition`, diagnostic errors are added using `diags.AddError` when decoding attributes, but the diagnostic is not checked or returned properly. This can lead to partially incorrect data being created if an error occurs, as the function will proceed and return an incomplete or default structure.

## Impact

High severity. Errors encountered during data marshalling or transformation are not propagated or handled meaningfully. This could lead to incorrect or incomplete outputs, and can cause subtle bugs which are difficult to track down during usage.

## Location

Lines 128–150 (example from `convertToDlpConnectorGroup`):

```go
func convertToDlpConnectorGroup(ctx context.Context, diags diag.Diagnostics, classification string, connectorsAttr basetypes.SetValue) dlpConnectorGroupsModelDto {
	var connectors []dataLossPreventionPolicyResourceConnectorModel
	err := connectorsAttr.ElementsAs(ctx, &connectors, true)
	if err != nil {
		diags.AddError("Client error when converting DlpConnectorGroups", "")
	}

	connectorGroup := dlpConnectorGroupsModelDto{
		Classification: classification,
		Connectors:     make([]dlpConnectorModelDto, 0),
	}

	for _, connector := range connectors {
		defaultAction := "Allow"

		if connector.DefaultActionRuleBehavior.ValueString() != "" {
			defaultAction = connector.DefaultActionRuleBehavior.ValueString()
		}

		connectorGroup.Connectors = append(connectorGroup.Connectors, dlpConnectorModelDto{
			Id:   connector.Id.ValueString(),
			Type: "Microsoft.PowerApps/apis",

			DefaultActionRuleBehavior: defaultAction,
			ActionRules:               convertToDlpActionRule(connector),
			EndpointRules:             convertToDlpEndpointRule(connector),
		})
	}
	return connectorGroup
}
```

Lines 163–175 (`convertToDlpCustomConnectorUrlPatternsDefinition`):

```go
func convertToDlpCustomConnectorUrlPatternsDefinition(ctx context.Context, diags diag.Diagnostics, connectorPatternsAttr basetypes.SetValue) []dlpConnectorUrlPatternsDefinitionDto {
	var customConnectorsPatterns []dataLossPreventionPolicyResourceCustomConnectorPattern
	err := connectorPatternsAttr.ElementsAs(ctx, &customConnectorsPatterns, true)
	if err != nil {
		diags.AddError("Client error when converting DlpCustomConnectorUrlPatternsDefinition", "")
	}

	customConnectorUrlPatternsDefinition := make([]dlpConnectorUrlPatternsDefinitionDto, 0)
	for _, customConnectorPattern := range customConnectorsPatterns {
		urlPattern := dlpConnectorUrlPatternsDefinitionDto{
			Rules: []dlpConnectorUrlPatternsRuleDto{},
		}
		urlPattern.Rules = append(urlPattern.Rules, dlpConnectorUrlPatternsRuleDto{
			Order:                       customConnectorPattern.Order.ValueInt64(),
			ConnectorRuleClassification: convertConnectorRuleClassificationValues(customConnectorPattern.DataGroup.ValueString()),
			Pattern:                     customConnectorPattern.HostUrlPattern.ValueString(),
		})
		customConnectorUrlPatternsDefinition = append(customConnectorUrlPatternsDefinition, urlPattern)
	}
	return customConnectorUrlPatternsDefinition
}
```

## Code Issue

```go
if err != nil {
	diags.AddError("Client error when converting DlpConnectorGroups", "")
}
...
return connectorGroup
```

```go
if err != nil {
	diags.AddError("Client error when converting DlpCustomConnectorUrlPatternsDefinition", "")
}
...
return customConnectorUrlPatternsDefinition
```

## Fix

Return early or propagate the error/diagnostics when an error is encountered. For example, you can return an error and handle it at a higher level, or ensure that the return values signal the error clearly.

Example for `convertToDlpConnectorGroup` (refactor return type to include error handling):

```go
func convertToDlpConnectorGroup(ctx context.Context, diags diag.Diagnostics, classification string, connectorsAttr basetypes.SetValue) (dlpConnectorGroupsModelDto, error) {
	var connectors []dataLossPreventionPolicyResourceConnectorModel
	err := connectorsAttr.ElementsAs(ctx, &connectors, true)
	if err != nil {
		diags.AddError("Client error when converting DlpConnectorGroups", err.Error())
		return dlpConnectorGroupsModelDto{}, err
	}
	// ...
	return connectorGroup, nil
}
```

Similarly, for `convertToDlpCustomConnectorUrlPatternsDefinition`:

```go
func convertToDlpCustomConnectorUrlPatternsDefinition(ctx context.Context, diags diag.Diagnostics, connectorPatternsAttr basetypes.SetValue) ([]dlpConnectorUrlPatternsDefinitionDto, error) {
	var customConnectorsPatterns []dataLossPreventionPolicyResourceCustomConnectorPattern
	err := connectorPatternsAttr.ElementsAs(ctx, &customConnectorsPatterns, true)
	if err != nil {
		diags.AddError("Client error when converting DlpCustomConnectorUrlPatternsDefinition", err.Error())
		return nil, err
	}
	// ...
	return customConnectorUrlPatternsDefinition, nil
}
```

