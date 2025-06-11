# Title

Mixed Naming Conventions for Field Names

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/models.go

## Problem

Go recommends using consistent naming conventions, typically CamelCase for struct fields and acronyms. In this file, `SubscriptionId` does not match Go's convention, where it should be `SubscriptionID`. Other fields such as `ID`, `AiType`, or `PackageName` are correct (ID/AI both in all-caps as is canonical for Go), but `SubscriptionId` uses a lowercase `d`.

## Impact

Severity: **Low**

This is mostly a readability and maintainability issue, but inconsistent naming can be confusing to contributors and breaks Go idioms, which may affect tooling and code generation further down the line.

## Location

Across all struct definitions in the file:

```go
	SubscriptionId    types.String `tfsdk:"subscription_id"`
```

## Fix

Rename the field to `SubscriptionID`. Ensure all code referencing this field is updated accordingly.

```go
	SubscriptionID    types.String `tfsdk:"subscription_id"`
```

Apply for the whole codebase
