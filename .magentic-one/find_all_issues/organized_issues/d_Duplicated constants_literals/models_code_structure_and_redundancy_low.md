# Code Structure and Redundancy

##

/workspaces/terraform-provider-power-platform/internal/services/connection/models.go

## Problem

There is noticeable repetition in the definition of data source and resource models for shares and connections. For example, `SharesDataSourceModel` and `SharesListDataSourceModel` follow a pattern that's duplicated by `ShareResourceModel`, `SharePrincipalResourceModel`, and their connection counterparts. These models could be parameterized or further composed to reduce boilerplate and the likelihood of inconsistencies as the codebase evolves.

## Impact

Low to Medium. Code repetition increases the risk of inconsistencies, makes updates more error-prone, and adds unnecessary maintenance overhead.

## Location

Widespread repetition in model definitions for shares and connections in:
- `SharesListDataSourceModel`
- `SharesDataSourceModel`
- `ConnectionsListDataSourceModel`
- `ConnectionsDataSourceModel`
- `ShareResourceModel`
- `SharePrincipalResourceModel`
- `ResourceModel`

## Code Issue

```go
// Example of pattern repetition
type SharesListDataSourceModel struct {
	Timeouts      timeouts.Value          `tfsdk:"timeouts"`
	EnvironmentId types.String            `tfsdk:"environment_id"`
	ConnectorName types.String            `tfsdk:"connector_name"`
	ConnectionId  types.String            `tfsdk:"connection_id"`
	Shares        []SharesDataSourceModel `tfsdk:"shares"`
}

type ConnectionsListDataSourceModel struct {
	Timeouts      timeouts.Value               `tfsdk:"timeouts"`
	EnvironmentId types.String                 `tfsdk:"environment_id"`
	Connections   []ConnectionsDataSourceModel `tfsdk:"connections"`
}
```

## Fix

Investigate whether generics, embedding, or composition can help reduce redundancy. For example, extract shared fields into base structs or use type embedding:

```go
type ListDataSourceModelBase struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
}

type SharesListDataSourceModel struct {
	ListDataSourceModelBase
	ConnectorName types.String            `tfsdk:"connector_name"`
	ConnectionId  types.String            `tfsdk:"connection_id"`
	Shares        []SharesDataSourceModel `tfsdk:"shares"`
}

type ConnectionsListDataSourceModel struct {
	ListDataSourceModelBase
	Connections []ConnectionsDataSourceModel `tfsdk:"connections"`
}
```

Adopting such refactoring will decrease maintenance effort and improve code clarity.
