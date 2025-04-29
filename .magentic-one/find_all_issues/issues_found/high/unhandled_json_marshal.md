# Title

Unhandled Errors from `json.Marshal`

##

`/internal/services/connection/datasource_connections.go`

## Problem

There are several instances in the code (inside the `ConvertFromConnectionDto` function) where errors from the `json.Marshal` function are ignored. Specifically, the following lines:

```go
if connection.Properties.ConnectionParametersSet != nil {
    p, _ := json.Marshal(connection.Properties.ConnectionParametersSet)
    conn.ConnectionParametersSet = types.StringValue(string(p))
}

if connection.Properties.ConnectionParameters != nil {
    p, _ := json.Marshal(connection.Properties.ConnectionParameters)
    conn.ConnectionParameters = types.StringValue(string(p))
}
```

In both cases, the returned error from `json.Marshal` is being ignored (`_` is used). This causes any potential marshalling failures to remain silent and unaddressed.

## Impact

- Ignoring errors could result in invalid or incomplete data being processed without any indications of issues, making it difficult to debug and troubleshoot.
- If an error occurs during marshalling, it could lead to inconsistencies or incorrect values being set in the `conn` struct.
- Severity: **High**

## Location

`ConvertFromConnectionDto` function, around lines:

```go
if connection.Properties.ConnectionParametersSet != nil {
    p, _ := json.Marshal(connection.Properties.ConnectionParametersSet)
}

if connection.Properties.ConnectionParameters != nil {
    p, _ := json.Marshal(connection.Properties.ConnectionParameters)
```

## Code Issue

```go
if connection.Properties.ConnectionParametersSet != nil {
    p, _ := json.Marshal(connection.Properties.ConnectionParametersSet)
    conn.ConnectionParametersSet = types.StringValue(string(p))
}

if connection.Properties.ConnectionParameters != nil {
    p, _ := json.Marshal(connection.Properties.ConnectionParameters)
    conn.ConnectionParameters = types.StringValue(string(p))
}
```

## Fix

Handle the errors returned by `json.Marshal` appropriately. For example:

```go
if connection.Properties.ConnectionParametersSet != nil {
    p, err := json.Marshal(connection.Properties.ConnectionParametersSet)
    if err != nil {
        // Log or handle the error appropriately, perhaps return it from this function
        return ConnectionsDataSourceModel{}, fmt.Errorf("failed to marshal ConnectionParametersSet: %w", err)
    }
    conn.ConnectionParametersSet = types.StringValue(string(p))
}

if connection.Properties.ConnectionParameters != nil {
    p, err := json.Marshal(connection.Properties.ConnectionParameters)
    if err != nil {
        // Log or handle the error appropriately, perhaps return it from this function
        return ConnectionsDataSourceModel{}, fmt.Errorf("failed to marshal ConnectionParameters: %w", err)
    }
    conn.ConnectionParameters = types.StringValue(string(p))
}
```

Ensure that the errors are actively logged, returned, or addressed to avoid silent failures.