# Panic-prone slice index access

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

In `convertPolicyModelToDlpPolicy`, within this loop

```go
for _, policy := range policy.CustomConnectorUrlPatternsDefinition {
	policyToCreate.CustomConnectorUrlPatternsDefinition.Rules = append(policyToCreate.CustomConnectorUrlPatternsDefinition.Rules, dlpConnectorUrlPatternsRuleDto{
		Order:                       policy.Rules[0].Order,
		ConnectorRuleClassification: policy.Rules[0].ConnectorRuleClassification,
		Pattern:                     policy.Rules[0].Pattern,
	})
}
```

there is no check that `policy.Rules` has at least one entry before accessing `policy.Rules[0]`. If `policy.Rules` is empty, this will panic.

## Impact

Could crash the application at runtime if the source data is malformed or the slice is empty. (Severity: High)

## Location

Lines 114-119

## Code Issue

```go
for _, policy := range policy.CustomConnectorUrlPatternsDefinition {
	policyToCreate.CustomConnectorUrlPatternsDefinition.Rules = append(policyToCreate.CustomConnectorUrlPatternsDefinition.Rules, dlpConnectorUrlPatternsRuleDto{
		Order:                       policy.Rules[0].Order,
		ConnectorRuleClassification: policy.Rules[0].ConnectorRuleClassification,
		Pattern:                     policy.Rules[0].Pattern,
	})
}
```

## Fix

Add a length check before accessing the first element of `policy.Rules`.

```go
for _, urlPattern := range policy.CustomConnectorUrlPatternsDefinition {
	if len(urlPattern.Rules) > 0 {
		policyToCreate.CustomConnectorUrlPatternsDefinition.Rules = append(
			policyToCreate.CustomConnectorUrlPatternsDefinition.Rules,
			dlpConnectorUrlPatternsRuleDto{
				Order:                       urlPattern.Rules[0].Order,
				ConnectorRuleClassification: urlPattern.Rules[0].ConnectorRuleClassification,
				Pattern:                     urlPattern.Rules[0].Pattern,
			},
		)
	}
}
```

