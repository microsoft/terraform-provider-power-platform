# Title

Repeated Code Logic in Create, Update, and Read Functions for Connection Parameters Parsing

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection.go

## Problem

The `Create`, `Update`, and `Read` functions contain repeated code logic for parsing and converting connection parameters (`ConnectionParameters` and `ConnectionParametersSet`). This violates the DRY (Don't Repeat Yourself) principle and makes the code less maintainable.

## Impact

- **Severity:** High  
- Repeated code can lead to inconsistent behavior if changes are applied in one function but not in others.  
- Maintenance complexity increases with code duplication.  
- In case of bugs or updates, developers need to make changes in multiple places, increasing the likelihood of errors.

## Location

Multiple occurrences in the `Create`, `Update`, and `Read` methods.

### Example (Create Function)
```go
if !plan.ConnectionParameters.IsNull() && plan.ConnectionParameters.ValueString() != "" {
    var params map[string]any
    err := json.Unmarshal([]byte(plan.ConnectionParameters.ValueString()), &params)
    if err != nil {
        resp.Diagnostics.AddError("Failed to convert connection parameters", err.Error())
        return
    }
    connectionToCreate.Properties.ConnectionParameters = params
}

if !plan.ConnectionParametersSet.IsNull() && plan.ConnectionParametersSet.ValueString() != "" {
    var params map[string]any
    err := json.Unmarshal([]byte(plan.ConnectionParametersSet.ValueString()), &params)
    if err != nil {
        resp.Diagnostics.AddError("Failed to convert connection parameters set", err.Error())
        return
    }
    connectionToCreate.Properties.ConnectionParametersSet = params
}
```

### Example (Update Function)
```go
var connParams map[string]any
if !plan.ConnectionParameters.IsNull() && plan.ConnectionParameters.ValueString() != "" {
    err := json.Unmarshal([]byte(plan.ConnectionParameters.ValueString()), &connParams)
    if err != nil {
        resp.Diagnostics.AddError("Failed to convert connection parameters", err.Error())
        return
    }
}

var connParamsSet map[string]any
if !plan.ConnectionParametersSet.IsNull() && plan.ConnectionParametersSet.ValueString() != "" {
    err := json.Unmarshal([]byte(plan.ConnectionParametersSet.ValueString()), &connParamsSet)
    if err != nil {
        resp.Diagnostics.AddError("Failed to convert connection parameters set", err.Error())
        return
    }
}
```

## Fix

### Suggested Refactor
Create a helper function to encapsulate the logic for parameter conversion and reuse it in all the functions.

```go
func parseConnectionParameters(parameters types.String) (map[string]any, error) {
    if parameters.IsNull() || parameters.ValueString() == "" {
        return nil, nil
    }
    var params map[string]any
    err := json.Unmarshal([]byte(parameters.ValueString()), &params)
    if err != nil {
        return nil, err
    }
    return params, nil
}

// Example usage in Create function
connParams, err := parseConnectionParameters(plan.ConnectionParameters)
if err != nil {
    resp.Diagnostics.AddError("Failed to convert connection parameters", err.Error())
    return
}
connectionToCreate.Properties.ConnectionParameters = connParams

connParamsSet, err := parseConnectionParameters(plan.ConnectionParametersSet)
if err != nil {
    resp.Diagnostics.AddError("Failed to convert connection parameters set", err.Error())
    return
}
connectionToCreate.Properties.ConnectionParametersSet = connParamsSet
```

Explanation: By creating a helper function, the code becomes more modular and easier to maintain. Updates to the conversion logic can be applied in a single place, reducing the risk of inconsistencies.