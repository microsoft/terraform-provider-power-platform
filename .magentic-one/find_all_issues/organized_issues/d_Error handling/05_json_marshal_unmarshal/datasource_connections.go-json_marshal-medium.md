# Issue: Unhandled errors when marshaling JSON in ConvertFromConnectionDto

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections.go

## Problem

In `ConvertFromConnectionDto`, calls to `json.Marshal` ignore returned errors, using the blank identifier `_`. If marshaling fails, this will silently set a possibly empty or wrong string value, introducing silent data loss or incorrect data.

## Impact

Severity: **Medium**

Silently ignoring marshaling errors can result in invalid or empty JSON strings in the state model, making debugging difficult and potentially leading to inaccurate resource states in Terraform.

## Location

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

## Code Issue

```go
p, _ := json.Marshal(connection.Properties.ConnectionParametersSet)
```

## Fix

Check the error returned by `json.Marshal` and handle it accordingly. At a minimum, you should potentially return a diagnostic or at least avoid setting the value if there's an error:

```go
if connection.Properties.ConnectionParametersSet != nil {
    p, err := json.Marshal(connection.Properties.ConnectionParametersSet)
    if err == nil {
        conn.ConnectionParametersSet = types.StringValue(string(p))
    } else {
        // Optionally log or handle the error here
    }
}

if connection.Properties.ConnectionParameters != nil {
    p, err := json.Marshal(connection.Properties.ConnectionParameters)
    if err == nil {
        conn.ConnectionParameters = types.StringValue(string(p))
    } else {
        // Optionally log or handle the error here
    }
}
```
