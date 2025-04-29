# Title

Inconsistent use of field names across types

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/dto.go`

## Problem

There is an inconsistency in naming conventions for similar types. For example:
- `DefaultConnectorActionRuleBehavior` in `dlpConnectorActionConfigurationsDto` uses PascalCase, while other fields follow CamelCase. 

This inconsistency can confuse contributors and cause issues, particularly for code readability and integration.

## Impact

Lack of standardization may lead to:
- Maintenance challenges.
- Misalignment with JSON field names when serialized.
- Increased risk of introducing errors when using the types across the codebase.

Severity: **Medium**

## Location

`dlpConnectorActionConfigurationsDto` struct type in `/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/dto.go`

## Code Issue

```go
type dlpConnectorActionConfigurationsDto struct {
	ConnectorId                        string             `json:"connectorId"`
	DefaultConnectorActionRuleBehavior string             `json:"defaultConnectorActionRuleBehavior"`
	ActionRules                        []dlpActionRuleDto `json:"actionRules"`
}
```

## Fix

Rename the field `DefaultConnectorActionRuleBehavior` to use CamelCase for consistent naming as with other fields.

```go
type dlpConnectorActionConfigurationsDto struct {
	ConnectorId                 string             `json:"connectorId"`
	DefaultActionRuleBehavior   string             `json:"defaultConnectorActionRuleBehavior"`
	ActionRules                 []dlpActionRuleDto `json:"actionRules"`
}
```