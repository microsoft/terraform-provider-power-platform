# Functions and variable naming consistency

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

There are some naming inconsistencies and potential confusion in function and variable names, such as `ConvertFromConnectionSharesDto` (using PascalCase for a converting function, and mixing DTO and model names), `SharesDataSourceModel`/`SharesListDataSourceModel` (potentially redundant or confusing), and the difference between singular/plural (`share` vs `shares`).

## Impact

**Severity: low**

This may lead to confusion for future contributors, reduce maintainability, and increase the risk of subtle bugs or misunderstandings in usage.

## Location

- `ConvertFromConnectionSharesDto`
- `SharesDataSourceModel` vs `SharesListDataSourceModel`
- `NewConnectionSharesDataSource`

## Code Issue

```go
func ConvertFromConnectionSharesDto(connection shareConnectionResponseDto) SharesDataSourceModel
```

## Fix

Adopt consistent, idiomatic Go naming conventions:

- Use camelCase for local variables and PascalCase for exported types/functions.
- Stick to singular/plural conventions (e.g., `ShareDataSourceModel`, `SharesListDataSourceModel`).
- Use clearer function names, such as `convertConnectionShareDtoToModel`.

Example:

```go
func convertConnectionShareDTOToModel(dto shareConnectionResponseDto) ShareDataSourceModel
```
