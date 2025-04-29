### Title

Ineffective Error Handling in `Read` Method

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/datasource_environment_templates.go

### Problem

Error handling in the `Read` function relies solely on a high-level error diagnostic for `GetEnvironmentTemplatesByLocation`. While this is functional, it does not contextualize the error, such as differentiating between network failures, improper authentication, or invalid parameters.

```go
environment_templates, err := d.EnvironmentTemplatesClient.GetEnvironmentTemplatesByLocation(ctx, state.Location.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}
```

### Impact

Without granular diagnostics, troubleshooting failures becomes challenging for users/devs, leading to delayed fixes and degraded user experience.

Severity: **Medium**

### Location

```go
environment_templates, err := d.EnvironmentTemplatesClient.GetEnvironmentTemplatesByLocation(ctx, state.Location.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}
```

### Fix

Refactor error handling to add detailed context based on the nature of the error:

```go
environment_templates, err := d.EnvironmentTemplatesClient.GetEnvironmentTemplatesByLocation(ctx, state.Location.ValueString())
if err != nil {
    if api.IsNetworkError(err) {
        resp.Diagnostics.AddError(
            "Network Error",
            fmt.Sprintf("Network error occurred while fetching templates: %s. Please verify your connection and try again.", err.Error()),
        )
    } else if api.IsAuthenticationError(err) {
        resp.Diagnostics.AddError(
            "Authentication Error",
            fmt.Sprintf("Authentication error occurred: %s. Please verify your credentials and try again.", err.Error()),
        )
    } else if api.IsInvalidParameterError(err) {
        resp.Diagnostics.AddError(
            "Invalid Parameters",
            fmt.Sprintf("Invalid parameters provided: %s. Please verify your configuration.", err.Error()),
        )
    } else {
        resp.Diagnostics.AddError(
            "Unknown Error",
            fmt.Sprintf("An unknown error occurred while fetching templates: %s. Please contact support.", err.Error()),
        )
    }
    return
}
```