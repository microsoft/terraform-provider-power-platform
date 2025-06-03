# SharesDataSource struct definition missing

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

The `SharesDataSource` struct is being used throughout the file, but there is no definition of this struct in the current file. If the struct is not declared elsewhere in the same package, this will lead to a compilation error. It is important for both readability and correctness that type declarations are visible or properly imported.

## Impact

**Severity: high**

The absence of the struct definition makes the code non-compilable. Functions like `Metadata`, `Schema`, `Configure`, `Read` are all expecting a receiver of type `*SharesDataSource`, but the type is not defined. This is a critical problem that prevents the codebase from building and being used.

## Location

The entire file usage of `SharesDataSource`:

```go
func NewConnectionSharesDataSource() datasource.DataSource {
	return &SharesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "connection_shares",
		},
	}
}
...
func (d *SharesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
...
func (d *SharesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
...
func (d *SharesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
...
func (d *SharesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
...
```

## Code Issue

```go
// References to SharesDataSource are made, but its struct definition is missing.
```

## Fix

Define the `SharesDataSource` struct appropriately based on usage patterns observed in the file. It should embed or contain the relevant client and type info fields used.

```go
type SharesDataSource struct {
	TypeInfo          helpers.TypeInfo
	ProviderTypeName  string
	ConnectionsClient *connectionsClient // assuming this is the client wrapper used for API calls
}
```

Adjust the field types to match exactly what is needed/referenced throughout the implementation.
