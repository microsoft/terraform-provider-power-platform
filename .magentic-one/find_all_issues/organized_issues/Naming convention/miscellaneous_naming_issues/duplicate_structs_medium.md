# Redundant and Potentially Duplicate Structs

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/dto.go

## Problem

Several structs appear to have near-duplicate purposes/fields or differ only slightly in naming, casing, or constituent types. For example:

- `dlpConnectorGroupsModelDto` and `dlpConnectorGroupsDto`
- `dlpConnectorModelDto` and `dlpConnectorDto`

This redundancy causes confusion as to which type should be used when, increases maintenance burden, and risks errors in conversion between nearly identical types.

## Impact

Severity: **Medium**

- Increased maintenance complexity.
- Higher risk of bugs due to confusion over similar types.
- More code to update when business requirements change.

## Location

Examples:

```go
type dlpConnectorGroupsDto struct {
	Classification string            `json:"classification"`
	Connectors     []dlpConnectorDto `json:"connectors"`
}

type dlpConnectorGroupsModelDto struct {
	Classification string                 `json:"classification"`
	Connectors     []dlpConnectorModelDto `json:"connectors"`
}
```

and

```go
type dlpConnectorModelDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	DefaultActionRuleBehavior string
	ActionRules               []dlpActionRuleDto
	EndpointRules             []dlpEndpointRuleDto
}

type dlpConnectorDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
```

## Fix

Consolidate these types where possible. Use a single, well-defined struct for each unique concept, augmenting with optional fields as needed for cases with additional data:

```go
type DlpConnectorDto struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	DefaultActionRuleBehavior string `json:"defaultActionRuleBehavior,omitempty"`
	ActionRules []dlpActionRuleDto  `json:"actionRules,omitempty"`
	EndpointRules []dlpEndpointRuleDto `json:"endpointRules,omitempty"`
}
```

Update usages throughout the codebase to use the consolidated types.

