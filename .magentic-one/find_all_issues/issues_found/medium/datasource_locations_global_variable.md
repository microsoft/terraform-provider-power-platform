# Title

Over-reliance on Global Variable in `NewLocationsDataSource`

## 

`/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations.go`

## Problem

The `NewLocationsDataSource` function sets a global variable (`helpers.TypeInfo`) which is then used to infer the type information. This creates implicit dependencies on the global scope and reduces modularity.

## Impact

Relying on global variables can lead to unintended side effects, especially when conditions change in the application execution context. It makes debugging harder and increases coupling between components. Severity is rated as **Medium**, as it directly affects maintainability and modularity.

## Location

```go
func NewLocationsDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "locations",
		},
	}
}
```

## Fix

Instead of relying on global variables, explicitly pass type information as an argument to the constructor. This makes the function more modular and easier to test.

```go
func NewLocationsDataSource(typeName string) datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: typeName,
		},
	}
}
```

And update calls to this function to pass the necessary type name explicitly:

```go
NewLocationsDataSource("locations")
```