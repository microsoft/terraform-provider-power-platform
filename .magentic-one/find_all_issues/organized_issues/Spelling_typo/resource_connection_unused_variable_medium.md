# Title
Incorrect variable naming: `conectionState` typo (should be `connectionState`)

##
/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection.go

## Problem
There is a typo in the variable name `conectionState` (should be `connectionState`) throughout the Create, Read, and Update methods.

## Impact
Medium severity: Typos in variable names can lead to confusion and reduce code readability and maintainability. This doesn't break functionality but is a code quality concern.

## Location
Multiple occurrences in methods: Create, Read, Update:

## Code Issue
```go
conectionState := ConvertFromConnectionDto(*connection)
plan.Id = types.String(conectionState.Id)
// ... and similar lines in Read and Update
```

## Fix
Rename all occurrences of `conectionState` to `connectionState` for clarity and consistency.

```go
connectionState := ConvertFromConnectionDto(*connection)
plan.Id = types.String(connectionState.Id)
// ... and similar fixes in Read and Update
```

