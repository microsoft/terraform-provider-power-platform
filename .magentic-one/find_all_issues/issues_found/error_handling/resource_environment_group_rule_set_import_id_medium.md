# Title

Error handling: misalignment between import and resource attribute ID

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set.go

## Problem

During the import function, resource `ImportStatePassthroughID` populates `"environment_group_id"`, but the main resource schema marks `"id"` as computed. Inconsistent use of import path versus resource identity could potentially cause confusion or require additional mapping elsewhere.

## Impact

Medium.

- Could cause inconsistent import behavior, requiring manual state editing or documentation caveats.

## Location

```go
resource.ImportStatePassthroughID(ctx, path.Root("environment_group_id"), req, resp)
```

## Fix

Ensure resource schema `id` and import attribute match or add import logic to map them correctly, depending on provider and Terraform expectations. Consider mapping ID property as a computed and/or set attribute.
