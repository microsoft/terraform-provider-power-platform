# Title

Hardcoded API Versions in Code

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_groups/api_environment_group.go`

## Problem

Throughout the file, API version strings are hardcoded directly in multiple places. For example, `"2021-04-01"` is used directly as the API version. Hardcoding like this makes future updates and maintenance harder, as changing the API version requires editing multiple places manually instead of updating a single constant or configuration.

## Impact

- **Severity:** High  
- Hardcoded values reduce maintainability and increase the chances of bugs when updating API versions. This pattern might also lead to inconsistent usage if different versions are hardcoded in various parts of the code.

## Location

Multiple locations in the file:

### `CreateEnvironmentGroup`
```go
	values.Add("api-version", "2021-04-01")
```

### `DeleteEnvironmentGroup`
```go
	values.Add("api-version", "2021-04-01")
```

### `UpdateEnvironmentGroup`
```go
	values.Add("api-version", "2021-04-01")
```

### `GetEnvironmentGroup`
```go
	values.Add("api-version", "2021-04-01")
```

### `GetEnvironmentsInEnvironmentGroup`
```go
	values.Add("api-version", "2021-04-01")
```

### `RemoveEnvironmentFromEnvironmentGroup`
```go
	values.Add("api-version", "1")
```

## Code Issue

```go
	values.Add("api-version", "2021-04-01")
	values.Add("api-version", "1")
```

## Fix

Define constants or utilize configuration management for API versioning.

```go
const ApiVersion2021 = "2021-04-01"
const ApiVersion1 = "1"

// Then use these constants in the code:
values.Add("api-version", ApiVersion2021)
values.Add("api-version", ApiVersion1)
```

Alternatively, they can be stored in configuration files or environment variables for dynamic management.
