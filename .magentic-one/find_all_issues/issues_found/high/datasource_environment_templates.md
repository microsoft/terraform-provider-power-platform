### Title

Potential Null Pointer Dereference in `Configure` Method

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/datasource_environment_templates.go

### Problem

The `Configure` method proceeds to use `client.Api` without adequate checks to ensure it's nil-safe. This can result in a nil pointer dereference if `client.Api` is nil:

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if ok {
    d.EnvironmentTemplatesClient = newEnvironmentTemplatesClient(client.Api)
}
```

Currently, the code assumes that `client.Api` will never be nil, which is not a guarded assumption.

### Impact

Dereferencing a nil pointer leads to a runtime panic and will cause the application to fail abruptly. This can affect user experience and necessitate provider troubleshooting in production environments.

Severity: **High**

### Location

```go
func (d *EnvironmentTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    client, ok := req.ProviderData.(*api.ProviderClient)
    if ok {
        d.EnvironmentTemplatesClient = newEnvironmentTemplatesClient(client.Api) // Potential nil pointer dereference.
    }
}
```

### Fix

Introduce a nil check for `client.Api` before proceeding to use it:

```go
func (d *EnvironmentTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    client, ok := req.ProviderData.(*api.ProviderClient)
    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected ProviderData Type",
            fmt.Sprintf("Expected *api.ProviderClient, but got: %T", req.ProviderData),
        )
        return
    }

    if client.Api == nil {
        resp.Diagnostics.AddError(
            "Invalid Client Api",
            "The Api object is nil in the provider client. This is an unexpected state. Please contact support or check the client configuration.",
        )
        return
    }

    d.EnvironmentTemplatesClient = newEnvironmentTemplatesClient(client.Api)
}
```