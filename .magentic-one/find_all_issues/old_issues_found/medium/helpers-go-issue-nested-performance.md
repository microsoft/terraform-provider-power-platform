# Title
Inefficient iteration logic in `convertToAttrValueCustomConnectorUrlPatternsDefinition`.

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem
The function uses nested loops unnecessarily inside the `convertToAttrValueCustomConnectorUrlPatternsDefinition`. For each `connectorUrlPattern`, the inner loop iterates through `Rules`, appending information into `connUrlPattern`. This logic fails to handle scalability when dealing with a potentially large dataset, degrading performance significantly.

---

## Impact
This can lead to degraded performance in systems with large `Rules` datasets. Overloading the system could eventually lead to timeouts or memory bottlenecks in extreme cases. This issue has **medium severity** as it impacts performance.

---

## Location
Line: 95

## Code Issue

```go
	for _, connectorUrlPattern := range urlPatterns {
		for _, rules := range connectorUrlPattern.Rules {
			connUrlPattern = append(connUrlPattern, types.ObjectValueMust(
				map[string]attr.Type{
					"order":            types.Int64Type,
					"host_url_pattern": types.StringType,
					"data_group":       types.StringType,
				},
				map[string]attr.Value{
					"order":            types.Int64Value(rules.Order),
					"host_url_pattern": types.StringValue(rules.Pattern),
					"data_group":       types.StringValue(convertConnectorRuleClassificationValues(rules.ConnectorRuleClassification)),
				},
			))
		}
	}
```
---

## Fix

```go
// Combine/Optimize appending logic to avoid deeply nested loops.
for _, connectorUrlPattern := range urlPatterns {
	ruleValues := make([]attr.Value, len(connectorUrlPattern.Rules))
	for i, rules := range connectorUrlPattern.Rules {
		ruleValues[i] = types.ObjectValueMust(
			map[string]attr.Type{
				"order":            types.Int64Type,
				"host_url_pattern": types.StringType,
				"data_group":       types.StringType,
			},
			map[string]attr.Value{
				"order":            types.Int64Value(rules.Order),
				"host_url_pattern": types.StringValue(rules.Pattern),
				"data_group":       types.StringValue(convertConnectorRuleClassificationValues(rules.ConnectorRuleClassification)),
			},
		)
	}
	connUrlPattern = append(connUrlPattern, ruleValues...)
}
```

Explanation: Adjusted logic to pre-calculate subset data for loops prior to appending into `connUrlPattern`.