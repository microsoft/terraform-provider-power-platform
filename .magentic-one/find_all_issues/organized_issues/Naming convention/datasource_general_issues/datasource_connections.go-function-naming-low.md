# Issue: Function Naming - `ConvertFromConnectionDto` Is Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections.go

## Problem

The function `ConvertFromConnectionDto` does not follow Go idioms for naming conversion methods. Go typically uses the form `toX` or `fromX` for conversion helpers, e.g., `connectionDtoToModel` or `toConnectionsDataSourceModel`.

## Impact

Severity: **Low**

While not a functional problem, non-idiomatic naming can make the codebase less consistent and harder to navigate for Go developers, especially in larger codebases.

## Location

```go
func ConvertFromConnectionDto(connection connectionDto) ConnectionsDataSourceModel
```

## Code Issue

```go
func ConvertFromConnectionDto(connection connectionDto) ConnectionsDataSourceModel
```

## Fix

Rename the function to follow Go conventions, such as:

```go
func connectionDtoToModel(connection connectionDto) ConnectionsDataSourceModel
```
