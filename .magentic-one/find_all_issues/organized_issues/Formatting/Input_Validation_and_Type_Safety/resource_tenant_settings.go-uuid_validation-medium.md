# Lack of input validation for critical fields

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

Some string attributes that represent UUID values (such as `environment_routing_target_environment_group_id` and `environment_routing_target_security_group_id`) do not have format validators to ensure valid UUIDs are supplied.

## Impact

Invalid data can be accepted at plan/apply time, and may lead to run-time errors or unexpected behavior interacting with the API. Data consistency and user experience are negatively affected. Severity: medium.

## Location

Schema definition for fields with `CustomType: customtypes.UUIDType{}` (e.g., under `.governance` nested field).

## Code Issue

```go
"environment_routing_target_environment_group_id": schema.StringAttribute{
    MarkdownDescription: "Assign newly created personal developer environments to a specific environment group",
    Optional:            true,
    CustomType:          customtypes.UUIDType{},
},
"environment_routing_target_security_group_id": schema.StringAttribute{
    MarkdownDescription: "Restrict routing to members of the following security group. (00000000-0000-0000-0000-000000000000 allows all users)",
    Optional:            true,
    CustomType:          customtypes.UUIDType{},
},
```

## Fix

Add UUID format validators. For example:

```go
import "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

...

"environment_routing_target_environment_group_id": schema.StringAttribute{
    MarkdownDescription: "Assign newly created personal developer environments to a specific environment group",
    Optional:            true,
    CustomType:          customtypes.UUIDType{},
    Validators: []validator.String{
        stringvalidator.RegexMatches(
            regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"),
            "Must be a valid UUID",
        ),
    },
},
```

And similarly for other UUID fields.

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_tenant_settings.go-uuid_validation-medium.md`
