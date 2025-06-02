# Title

Potential Null Dereference in `convertPolicyModelToDlpPolicy`

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

In `convertPolicyModelToDlpPolicy`, the conditional check for `policy.ConnectorConfigurationsDefinition != nil` is insufficient, as subsequent assignment operations might cause a null reference exception for ConnectorConfigurationsDefinition in edge cases.

## Impact

Can lead to runtime errors in cases where the `ConnectorConfigurationsDefinition` is unintentionally null and operations are performed on the null object.

Severity marked **High**.

## Location

`func convertPolicyModelToDlpPolicy(policy dlpPolicyModelDto) dlpPolicyDto`

## Code Issue

```go
policyToCreate.ConnectorConfigurationsDefinition = nil
if len(connectorActionConfigurationsDto) > 0 || len(endpointConfigurationsDto) > 0 {
policyToCreate.ConnectorConfigurationsDefinition = &dlpConnectorConfigurationsDefinitionDto{}

if len(connectorActionConfigurationsDto) > 0 {
policyToCreate.ConnectorConfigurationsDefinition.ConnectorActionConfigurations = connectorActionConfigurationsDto
}
if len(endpointConfigurationsDto) > 0 {
policyToCreate.ConnectorConfigurationsDefinition.EndpointConfigurations = endpointConfigurationsDto
}
} else {
policyToCreate.ConnectorConfigurationsDefinition = nil
}
```

## Fix

Ensure null checks and use proper initialization.

```go
if len(connectorActionConfigurationsDto) > 0 || len(endpointConfigurationsDto) > 0 {
policyToCreate.ConnectorConfigurationsDefinition = &dlpConnectorConfigurationsDefinitionDto{
ConnectorActionConfigurations: []dlpConnectorActionConfigurationsDto{},
EndpointConfigurations: []dlpEndpointConfigurationsDto{},
}
if len(connectorActionConfigurationsDto) > 0 {
policyToCreate.ConnectorConfigurationsDefinition.ConnectorActionConfigurations = connectorActionConfigurationsDto
}
if len(endpointConfigurationsDto) > 0 {
policyToCreate.ConnectorConfigurationsDefinition.EndpointConfigurations = endpointConfigurationsDto
}
} else {
policyToCreate.ConnectorConfigurationsDefinition = nil
}
```