# Title

Improper Error Handling When Converting Columns to Map

## Path

`/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go`

## Problem

When invoking the `convertResourceModelToMap` function in multiple methods (`Create`, `Delete`, `Update`), the returned error is not sufficiently validated or investigated. Specifically, in cases like this:

```go
mapColumns, err := convertResourceModelToMap(&stateColumns)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Error converting columns to map: %s", err.Error()), err.Error())
    return
}
```

The absence of context or diagnostics on why the `convertResourceModelToMap` function fails could hinder effective debugging.

## Impact

This creates insufficient visibility into root causes during runtime, impacting the diagnosability and maintainability of the system. Severity: **Critical**

## Location

- `Create` function
- `Delete` function
- `Update` function
- `Read` function (partially related)

## Code Issue

```go
mapColumns, err := convertResourceModelToMap(&stateColumns)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Error converting columns to map: %s", err.Error()), err.Error())
    return
}
```

## Fix

Enhance error logging and context by including more information about method variables and the specific state leading to the error. For example:

```go
mapColumns, err := convertResourceModelToMap(&stateColumns)
if err != nil {
    resp.Diagnostics.AddError(
        "Error converting columns to map",
        fmt.Sprintf("Error occurred while processing stateColumns: %v. Full error: %v", stateColumns, err.Error()),
    )
    return
}
```