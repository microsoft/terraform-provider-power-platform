# Title

Usage of `stringplanmodifier.RequiresReplace()` in `Schema` for `id` attribute

## Path to the file

`/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go`

## Problem

The `id` attribute uses `stringplanmodifier.RequiresReplace()` in its `PlanModifiers` configuration. While this enforces replacement behavior when the value changes, this could unintentionally trigger resource recreation even when only minor updates or corrections requiring less disruptive changes occur. This might be overly aggressive for certain scenarios.

## Impact

The aggressive behavior of `RequiresReplace()` can lead to unnecessary recreation of resources, which may result in downtime or unintended side effects for dependent resources. Severity level is **medium** since the impact is undesirable but not fundamentally destructive.

## Location

```go
"id": schema.StringAttribute{
  MarkdownDescription: "Client id for the service principal",
  Required:            true,
  CustomType:          customtypes.UUIDType{},
  PlanModifiers: []planmodifier.String{
      stringplanmodifier.RequiresReplace(),
  },
},
```

## Fix

To make the replacement behavior less aggressive, consider using a custom plan modifier or revising your resource architecture to minimize unnecessary recreation. Utilize custom validation mechanisms to provide greater control:

### Suggested Fix:

```go
"id": schema.StringAttribute{
    MarkdownDescription: "Client id for the service principal",
    Required:            true,
    CustomType:          customtypes.UUIDType{},
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.OptionalTriggerReplace(),
        // OR implement more refined custom logic
    },
},
```
