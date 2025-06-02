# Title

Duplicate Struct Definitions Leading to Redundancy

##

/workspaces/terraform-provider-power-platform/internal/services/solution/models.go

## Problem

The file defines both `DataSource` and `Resource` structs, each with identical fields (`helpers.TypeInfo`, `SolutionClient Client`). This duplication increases maintenance work and risk of divergence between logically similar structures.

## Impact

Medium. Structural duplication may lead to unintentional divergence of similar logic and additional maintenance overhead with every code or logic update involving these constructors.

## Location

- Definition of `DataSource` and `Resource` structs:

## Code Issue

```go
type DataSource struct {
	helpers.TypeInfo
	SolutionClient Client
}

type Resource struct {
	helpers.TypeInfo
	SolutionClient Client
}
```

## Fix

Combine the two structs if their logical use is indeed the same or inherit via embedding, or explain/document if they are meant to diverge in the future.

```go
// Example: Create a single struct to unify logic
type SolutionService struct {
	helpers.TypeInfo
	SolutionClient Client
}

// Or, if separation is required, provide explicit documentation for why both exist.
```
