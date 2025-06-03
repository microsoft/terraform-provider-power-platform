# Direct API Usage Without Typed Abstraction

##

/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations.go

## Problem

In the `Configure` method, the assignment:
```go
d.LocationsClient = newLocationsClient(client.Api)
```
suggests that the API structure (`client.Api`) is directly passed into a constructor function to create the internal `LocationsClient`. This could indicate an abstraction leak, where lower-level API details are exposed into the service/data source layer. This approach can hinder testing/mocking, tightly couple layers, and reduce extensibility.

## Impact

Direct propagation of lower-level API references increases long-term maintenance and refactoring risk. Changes in API details could ripple through consumer logic. **Severity: Medium**

## Location

In `Configure`:

```go
client, ok := req.ProviderData.(*api.ProviderClient)
// ...
d.LocationsClient = newLocationsClient(client.Api)
```

## Code Issue

```go
d.LocationsClient = newLocationsClient(client.Api)
```

## Fix

Hide low-level API details behind a well-defined interface or struct. Pass only necessary, well-typed dependencies. For example, if the intention is to mock or abstract the transport, define an interface on the API, and inject that.

```go
// Prefer some interface, e.g.,
type LocationsAPI interface {
    GetLocations(ctx context.Context) ([]Location, error)
}

d.LocationsClient = api.NewLocationsClient(client)
// Now LocationsClient is based on an interface, not concrete Api struct internals.
```

This approach increases maintainability, testing, and cleanliness.
