# Title

Missing Error Handling and Diagnostics in `PlanModifyBool`

## 

`/workspaces/terraform-provider-power-platform/internal/modifiers/restore_original_value_modifier.go`

## Problem

In the `PlanModifyBool` method, the code uses `resp.Private.SetKey` to store private data `[ ]byte{}` but lacks proper error handling and diagnostics for failed operations. Furthermore, there is no validation of serialization or context around what `[]byte{}` represents, leading to ambiguous and potentially unsafe operations.

## Impact

The lack of error handling reduces the reliability of the method because if the `SetKey` operation fails, this will go unnoticed, leaving the system in an inconsistent or incorrect state. It also increases the likelihood of issues in production. Severity: **High**.

## Location

File location: `/workspaces/terraform-provider-power-platform/internal/modifiers/restore_original_value_modifier.go`, within the `PlanModifyBool` function.

## Code Issue

```go
resp.Private.SetKey(ctx, req.Path.String(), []byte{})
```

## Fix

Introduce error handling and use a meaningful serialization mechanism for the stored data. The fix also includes adding diagnostics to report any error effectively:

```go
// Serialize meaningful data to store the original boolean value
boolData := map[string]interface{}{
    "originalValue": req.ConfigValue.Value(),
}

// Convert the data to byte format safely
serializedBoolData, err := json.Marshal(boolData)
if err != nil {
    resp.Diagnostics.AddError(
        "Failed to encode original boolean value",
        fmt.Sprintf("Error serializing original value for attribute %s: %s", req.PathExpression.String(), err.Error()),
    )
    return
}

// Safely store the key in private data with error handling
if err := resp.Private.SetKey(ctx, req.Path.String(), serializedBoolData); err != nil {
    resp.Diagnostics.AddError(
        "Failed to store original boolean value",
        fmt.Sprintf("Error storing original value for attribute %s: %s", req.PathExpression.String(), err.Error()),
    )
    return
}
```

This ensures safety and clarity in the implementation of the method.
