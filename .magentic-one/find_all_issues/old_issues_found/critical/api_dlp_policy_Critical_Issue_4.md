# Title

Cross-Iteration Redeclaration of Variables in `convertPolicyModelToDlpPolicy`

## Path

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

The method `convertPolicyModelToDlpPolicy`, when iterating through certain collections, repeatedly redeclares variable names within each iteration, leading to poor readability and increased risk of errors.

## Impact

This leads to readability issues, bloated logic, and potential hidden issues during large-scale maintenance. Severity marked **Critical** as such an approach is unsafe and misleading for developers.

## Location

`func convertPolicyModelToDlpPolicy(policy dlpPolicyModelDto) dlpPolicyDto`

## Code Issue

```go
for _, policy := range policy.CustomConnectorUrlPatternsDefinition {
policyToCreate.CustomConnectorUrlPatternsDefinition.Rules = append(policyToCreate.CustomConnectorUrlPatternsDefinition.Rules, dlpConnectorUrlPatternsRuleDto{
Order: policy.Rules[0].Order,
ConnectorRuleClassification: policy.Rules[0].ConnectorRuleClassification,
Pattern: policy.Rules[0].Pattern,
})
}

for _, connGroups := range policy.ConnectorGroups {
conG := dlpConnectorGroupsDto{
Classification: connGroups.Classification,
Connectors: []dlpConnectorDto{},
}

for _, connector := range connGroups.Connectors {
nameSplit := strings.Split(connector.Id, "/")
con := dlpConnectorDto{
Id: connector.Id,
Name: nameSplit[len(nameSplit)-1],
Type: connector.Type,
}
conG.Connectors = append(conG.Connectors, con)
}
policyToCreate.PolicyDefinition.ConnectorGroups = append(policyToCreate.PolicyDefinition.ConnectorGroups, conG)
}
```

## Fix

Focus on refactoring for clarity and uniqueness of iteration variable names and logic improvements.

```go
for _, urlPattern := range policy.CustomConnectorUrlPatternsDefinition {
policyToCreate.CustomConnectorUrlPatternsDefinition.Rules = append(policyToCreate.CustomConnectorUrlPatternsDefinition.Rules, dlpConnectorUrlPatternsRuleDto{
Order: urlPattern.Rules[0].Order,
ConnectorRuleClassification: urlPattern.Rules[0].ConnectorRuleClassification,
Pattern: urlPattern.Rules[0].Pattern,
})
}

for _, group := range policy.ConnectorGroups {
connectorGroup := dlpConnectorGroupsDto{
Classification: group.Classification,
Connectors: []dlpConnectorDto{},
}

for _, conn := range group.Connectors {
nameSplit := strings.Split(conn.Id, "/")
connector := dlpConnectorDto{
Id: conn.Id,
Name: nameSplit[len(nameSplit)-1],
Type: conn.Type,
}
connectorGroup.Connectors = append(connectorGroup.Connectors, connector)
}
policyToCreate.PolicyDefinition.ConnectorGroups = append(policyToCreate.PolicyDefinition.ConnectorGroups, connectorGroup)
}
```