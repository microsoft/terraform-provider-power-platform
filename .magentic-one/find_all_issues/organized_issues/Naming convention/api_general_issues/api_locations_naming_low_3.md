# Title

Return type `locationDto` does not follow Go type naming conventions

##

/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go

## Problem

Go type naming conventions recommend using PascalCase. The type `locationDto` should be renamed to `LocationDTO` or just `Location` if "DTO" is superfluous. This improves consistency throughout Go codebases.

## Impact

Low severity. Naming inconsistencies can reduce code readability and teamwork efficiency.

## Location

Return values in method `GetLocations`:

## Code Issue

```go
func (client *client) GetLocations(ctx context.Context) (locationDto, error) {
	// ...
	var locations locationDto
	// ...
}
```

## Fix

Rename `locationDto` to `LocationDTO` (or simpler, `Location`), and ensure the type declaration elsewhere follows the same.

```go
func (client *client) GetLocations(ctx context.Context) (LocationDTO, error) {
	// ...
	var locations LocationDTO
	// ...
}
```
