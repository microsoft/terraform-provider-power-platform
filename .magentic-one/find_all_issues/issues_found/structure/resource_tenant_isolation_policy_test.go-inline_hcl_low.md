# Title

Overuse of Inline Terraform HCL Strings Reduces Test Maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy_test.go

## Problem

Many test cases have Terraform configuration directly inlined as large, multi-line strings inside test steps. While this is common in provider tests, its repeated use for similar resources with only minor modifications can reduce readability and maintainability—especially as scenarios increase.

## Impact

When test resource configuration is duplicated this way, it becomes error-prone to update all copies if the underlying resource schema changes. It is also harder to visually compare what is changing from step to step, slowing test review and modification. This is a **low** severity issue, as correctness is not directly affected, but maintenance cost grows as the suite expands.

## Location

For example, in every resource test step:

```go
Config: `
    resource "powerplatform_tenant_isolation_policy" "test" {
        is_disabled = false
        allowed_tenants = toset([
            {
                tenant_id = "11111111-1111-1111-1111-111111111111"
                inbound  = true
                outbound = true
            }
        ])
    }`,
```

## Code Issue

```go
{
    Config: `
    resource "powerplatform_tenant_isolation_policy" "test" {
        is_disabled = false
        allowed_tenants = toset([
            {
                tenant_id = "11111111-1111-1111-1111-111111111111"
                inbound  = true
                outbound = true
            }
        ])
    }`,
    Check: ...
}
```

## Fix

Extract the common configurations into standalone helper constants or functions to increase clarity and reuse, for example:

```go
const tenantIsolationPolicyBase = `
resource "powerplatform_tenant_isolation_policy" "test" {
    is_disabled = %v
    allowed_tenants = toset([
        {
            tenant_id = "%s"
            inbound = %v
            outbound = %v
        }
    ])
}`

// And in the test:
Config: fmt.Sprintf(tenantIsolationPolicyBase, false, "11111111-1111-1111-1111-111111111111", true, true),
```

Use parameterized helpers or templates for variants, so changes are made in one place and intent is clearer.
