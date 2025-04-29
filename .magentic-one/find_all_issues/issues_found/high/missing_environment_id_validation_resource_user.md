# Title

Missing validation for environment IDs

##

`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go`

## Problem

Environment IDs (`environment_id`) appear multiple times in the code but lack strong validation to ensure that they adhere to expected formats (e.g., string `GUID`). Ensuring environment IDs conform to standards prevents errors during operations like `Create`, `Read`, and `Delete`.

## Impact

This issue is **high severity** because incorrect or malformed environment IDs could lead to failures during interaction with downstream API calls, potentially compromising the database and external integrations.

## Location

Example from the `Create` method:

```go
hasEnvDataverse, err := r.UserClient.EnvironmentHasDataverse(ctx, plan.EnvironmentId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
    return
}
```

## Code Issue

### Example

```go
plan.EnvironmentId.ValueString()
```

and

```go
resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
```

are called without prior validation.

## Fix

Introduce a dedicated validation function or middleware that ensures `environment_id` adheres to expected formats (e.g., valid GUID). Modify the schema definition to explicitly validate.

### Corrected Code

#### Add validation middleware
```go
func validateEnvironmentID(environmentID string) error {
    if !isValidGUID(environmentID) {
        return fmt.Errorf("Environment ID %s is not a valid GUID", environmentID)
    }
    return nil
}
```

#### Usage in `Create` method
```go
err := validateEnvironmentID(plan.EnvironmentId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Invalid environment ID provided %s", plan.EnvironmentId.ValueString()), err.Error())
    return
}
hasEnvDataverse, err := r.UserClient.EnvironmentHasDataverse(ctx, plan.EnvironmentId.ValueString())
```

#### Schema modification
```go
"environment_id": schema.StringAttribute{
    MarkdownDescription: "Unique environment id (guid). Must adhere to valid GUID format.",
    Required:            true,
    Validators: []validator.String{
        stringvalidator.Regex(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`),
    },
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
},
```

### Explanation

This fix ensures that `environment_id` is validated before usage, preventing API calls with invalid IDs and reducing debugging complexity.
