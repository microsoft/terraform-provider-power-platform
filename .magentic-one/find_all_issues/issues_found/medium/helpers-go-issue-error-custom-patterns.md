# Title
Improper error handling in `convertToDlpCustomConnectorUrlPatternsDefinition`.

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem
The error message added in case `ElementsAs` fails does not include specific details about the encountered error, making debugging harder.

---

## Impact
This makes troubleshooting more time-consuming and reduces the visibility of the issue's root cause during runtime. Considered **medium severity**, as it affects developer experience and debugging but does not break functionality.

---

## Location
Line: 370

## Code Issue

```go
	err := connectorPatternsAttr.ElementsAs(ctx, &customConnectorsPatterns, true)
	if err != nil {
		diags.AddError("Client error when converting DlpCustomConnectorUrlPatternsDefinition", "")
	}
```
---

## Fix

```go
// Append actual error details to the diagnostic message to improve clarity.
	err := connectorPatternsAttr.ElementsAs(ctx, &customConnectorsPatterns, true)
	if err != nil {
		diags.AddError("Error converting DlpCustomConnectorUrlPatternsDefinition",
		fmt.Sprintf("Details: %v", err))
	}
```

Explanation: Enhancing error messaging by feeding the actual error details into the diagnostic system.