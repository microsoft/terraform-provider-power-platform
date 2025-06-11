# Lack of State Initialization Before Read in Read Method

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors.go

## Problem

In the `Read` method, the `state` object is decoded from the request’s state, but there is no initialization to ensure `state.Connectors` is a non-nil (empty) slice prior to appending connector models. If `state.Connectors` is nil, repeated Read invocations or certain state transitions could result in a nil slice being marshalled to Terraform, which may not be handled consistently by the framework or downstream consumers.

## Impact

**Low to Medium severity:** Possible risk of data inconsistencies or unexpected nil slices propagating to the Terraform state, which could lead to intermittent deserialization issues or subtle schema mismatches.

## Location

```go
var state ListDataSourceModel
resp.State.Get(ctx, &state)

for _, connector := range connectors {
    connectorModel := convertFromConnectorDto(connector)
    state.Connectors = append(state.Connectors, connectorModel)
}
```

## Fix

Ensure `state.Connectors` is initialized to an empty slice if nil, before appending new elements:

```go
var state ListDataSourceModel
resp.State.Get(ctx, &state)
if state.Connectors == nil {
    state.Connectors = []ConnectorModel{}
}

for _, connector := range connectors {
    connectorModel := convertFromConnectorDto(connector)
    state.Connectors = append(state.Connectors, connectorModel)
}
```
