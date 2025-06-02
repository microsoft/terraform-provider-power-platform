# Title
Unnecessary early return inside loop causes incorrect flow.

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem
In the function `convertToAttrValueConnectorsGroup`, the loop prematurely returns a `types.SetValueMust` when a classification condition matches. This prevents complete processing of the connectors group.

---

## Impact
The loop aborts further execution upon finding a match, which can skip processing of other items and compromise data correctness or compliance. This is considered **high severity** because it directly affects the functional correctness of the application.

---

## Location
Line: 77

## Code Issue

```go
	for _, conn := range connectorsGroup {
		if conn.Classification == classification {
			return types.SetValueMust(connectorSetObjectType, convertToAttrValueConnectors(conn, connectorValues))
		}
	}
	return types.SetValueMust(connectorSetObjectType, []attr.Value{})
```
---

## Fix

```go
// Avoid early return inside the loop; collect data first then decide the return.
connectorResult := []attr.Value{}
for _, conn := range connectorsGroup {
	if conn.Classification == classification {
		connectorResult = convertToAttrValueConnectors(conn, connectorValues)
		break
	}
}
return types.SetValueMust(connectorSetObjectType, connectorResult)
```

Explanation: The fix ensures the loop completes execution by storing valid results in a variable `connectorResult` and returning only when the loop is finished.