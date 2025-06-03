# Use of hard-coded string for magic UUIDs

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

The schema field for `environment_routing_target_security_group_id` uses a hard-coded magic UUID value `00000000-0000-0000-0000-000000000000` in the description, but there is no reference to or enforcement for this value via constants. All code interacting with this field must treat this value specially, but the only place it is mentioned is in a schema description. This weakens data consistency and can lead to logic spread across the codebase.

## Impact

Hard to maintain, error-prone documentation and code; the expectation for a special UUID value is insulated in schema description only. If this logic is required elsewhere, it risks divergence and bugs. Severity: low.

## Location

Schema for `environment_routing_target_security_group_id`:

```go
"environment_routing_target_security_group_id": schema.StringAttribute{
    MarkdownDescription: "Restrict routing to members of the following security group. (00000000-0000-0000-0000-000000000000 allows all users)",
    Optional:            true,
    CustomType:          customtypes.UUIDType{},
},
```

## Fix

Declare a constant (if not already in use elsewhere) like:

```go
const ALLOW_ALL_USERS_UUID = "00000000-0000-0000-0000-000000000000"
```

Reference it in code and in descriptions by interpolating or documenting that constant. Additionally, logic that relies on this value should use this constant for data comparison and assignment, not inline string literals.

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_tenant_settings.go-magic_uuid-low.md`
