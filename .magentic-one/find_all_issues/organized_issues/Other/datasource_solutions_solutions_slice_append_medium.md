# Control Flow: Appending to Solutions Slice Without Resetting in Read

##

/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions.go

## Problem

In the `Read` method, `state.Solutions` can be appended to without clearing any existing data within it. If `resp.State.Get` deserializes a non-empty list (left over from previous partial state or refreshes), the `state.Solutions` will accumulate items across invocations, leading to duplication or stale results.

## Impact

**Medium**. This could introduce confusing, incorrect, or duplicated state for multiple solution reads/refreshes, leading to subtle bugs or incorrect data displayed to users.

## Location

Lines 121-126, in the `Read` method:

## Code Issue

```go
for _, solution := range solutions {
	solutionModel := convertFromSolutionDto(solution)
	state.Solutions = append(state.Solutions, solutionModel)
}
```

## Fix

Reset `state.Solutions` before populating it anew:

```go
state.Solutions = nil
for _, solution := range solutions {
	solutionModel := convertFromSolutionDto(solution)
	state.Solutions = append(state.Solutions, solutionModel)
}
```

Or alternatively, pre-size the slice if the full list is available and performance is a concern.
