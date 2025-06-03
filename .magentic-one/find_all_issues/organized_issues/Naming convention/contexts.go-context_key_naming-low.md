# Issue: Consistency in Context Key Naming

##

/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go

## Problem

The context keys (`EXECUTION_CONTEXT_KEY`, `REQUEST_CONTEXT_KEY`, `TEST_CONTEXT_KEY`) are all defined as constants of type `ContextKey`. However, using string keys in a context is generally discouraged for cross-package usage (to avoid collisions). Defining and using `ContextKey` as a specific type is correct. However, the naming convention for context keys usually uses lowerCamelCase when used as values rather than constants, to indicate unexported/private usage.

## Impact

- **Impact:** Low  
  This is a minor style/readability topic and is not likely to cause a bug, but alignment with norms helps reduce surprises for other developers.

## Location

Definition of context keys:

## Code Issue

```go
const (
	EXECUTION_CONTEXT_KEY ContextKey = "executionContext"
	REQUEST_CONTEXT_KEY   ContextKey = "requestContext"
	TEST_CONTEXT_KEY      ContextKey = "testContext"
)
```

## Fix

Consider following Go conventions—private (unexported) variables/constants should be lowerCamelCase.

```go
const (
	executionContextKey ContextKey = "executionContext"
	requestContextKey   ContextKey = "requestContext"
	testContextKey      ContextKey = "testContext"
)
```

Update usages accordingly.
