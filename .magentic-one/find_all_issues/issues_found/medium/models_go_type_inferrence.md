# Title

Type Safety and Missing Validations in Struct

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/models.go

## Problem

Struct field types are defined directly, but there are no validations for required fields such as `Location` or `Name`. This can lead to incomplete or invalid data being passed into these structures, which causes downstream issues.

## Impact

The missing validation impacts the accuracy and robustness of the code. This issue has a "medium" severity, as improper input limits functionality or causes runtime exceptions without proper validation logic.

## Location

Definition and implementation of structs `EnvironmentTemplatesDataSourceModel` and `EnvironmentTemplatesDataModel`.

## Code Issue

```go

type EnvironmentTemplatesDataModel struct {
    Category                     string `tfsdk:"category"`
    ID                           string `tfsdk:"id"`
    Name                         string `tfsdk:"name"`
    DisplayName                  string `tfsdk:"display_name"`
    Location                     string `tfsdk:"location"`
    IsDisabled                   bool   `tfsdk:"is_disabled"`
    DisabledReasonCode           string `tfsdk:"disabled_reason_code"`
    DisabledReasonMessage        string `tfsdk:"disabled_reason_message"`
    IsCustomerEngagement         bool   `tfsdk:"is_customer_engagement"`
    IsSupportedForResetOperation bool   `tfsdk:"is_supported_for_reset_operation"`
}

```

## Fix

Introduce validation logic for fields such as `Location`, `Name`, etc., to ensure that invalid data cannot pass through.

For example:
```go

func validateEnvironmentTemplate(data EnvironmentTemplatesDataModel) error {
    if data.Location == "" {
        return errors.New("location cannot be empty")
    }
    if data.Name == "" {
        return errors.New("name cannot be empty")
    }
    return nil
}

```