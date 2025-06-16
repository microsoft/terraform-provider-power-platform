# Title

Inconsistent Resource Naming: Typo in NewEnterpisePolicyResource

##

/workspaces/terraform-provider-power-platform/internal/provider/provider.go

## Problem

There is a typo in the resource provider function `enterprise_policy.NewEnterpisePolicyResource()`. The word "Enterpise" should be "Enterprise".

## Impact

Incorrect naming could cause confusion and violates naming consistency. It may also lead to import errors or break code referencing this function elsewhere if fixed without care. Severity: **medium**.

## Location

In the Resources registration function, specifically this line:

```go
func() resource.Resource { return enterprise_policy.NewEnterpisePolicyResource() },
```

## Fix

Rename the function call to use "Enterprise" instead of "Enterpise":

```go
func() resource.Resource { return enterprise_policy.NewEnterprisePolicyResource() },
```

Be sure to fix the actual function name in the implementation file under the enterprise_policy package as well to maintain consistency.
