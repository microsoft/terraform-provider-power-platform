# Title

Misleading Documentation for `Schema` Method Description

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set.go`

## Problem

The description provided for the `Schema` method mentions an API limitation ("Known Issue: This resource only works with a user context and can not be used at this time with a service principal.") but does not clarify under what specific conditions this occurs or whether there are alternative solutions. 

```go
MarkdownDescription: "Allows the creation of environment group rulesets. See [Power Platform documentation](https://learn.microsoft.com/power-platform/admin/environment-groups) for more information on the available rules that can be applied to an environment group.\n\n!> Known Issue: This resource only works with a user context and can not be used at this time with a service principal.  This is a limitation of the underlying API.",
```

## Impact

- Leads to user confusion, as the specifics of the API limitation and its workaround are not clear.
- Users may try to implement methods without understanding the limitation, wasting time and effort.
- Reduces usability of the resource due to poor documentation.

**Severity:** Low

## Location

Method Name: `Schema`

## Code Issue

```go
MarkdownDescription: "Allows the creation of environment group rulesets. See [Power Platform documentation](https://learn.microsoft.com/power-platform/admin/environment-groups) for more information on the available rules that can be applied to an environment group.\n\n!> Known Issue: This resource only works with a user context and can not be used at this time with a service principal.  This is a limitation of the underlying API.",
```

## Fix

Update the `MarkdownDescription` to clarify the issue and provide alternative solutions, if any.

```go
MarkdownDescription: "Allows the creation of environment group rulesets. See [Power Platform documentation](https://learn.microsoft.com/power-platform/admin/environment-groups) for more information on the available rules that can be applied to an environment group.\n\n!> Known Issue: This resource currently only supports user context operations and cannot be used with a service principal. This limitation is due to the underlying API restrictions. Users may need to ensure proper user authentication when interacting with this resource until further updates to the API address this limitation.",
```

- The updated note clarifies the limitations and offers actionable insights for users to mitigate the issue.