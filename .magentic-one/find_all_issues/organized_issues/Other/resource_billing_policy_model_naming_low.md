# Issue: Undocumented or Ambiguous Resource/Field Naming and Usage

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go

## Problem

The code uses struct names such as `BillingPolicyResourceModel`, `billingPolicyCreateDto`, and `BillingPolicyUpdateDto` but these types are not defined in the given file. This potentially leads to confusion: a reader cannot understand the data model, required and computed fields without searching elsewhere. This harms code maintainability and can make onboarding contributors more difficult.

## Impact

Severity: **Low**

Lack of clear definitions or comments for implicit or externally referenced types reduces code readability, increases context-switching costs for contributors, and may cause developers to reimplement logic that already exists elsewhere.

## Location

- Resource methods, plan/state population, and API interactions:
  ```go
  var plan *BillingPolicyResourceModel
  billingPolicyToCreate := billingPolicyCreateDto{...}
  ```
- Nowhere in this file are these types defined or described.

## Code Issue

```go
var plan *BillingPolicyResourceModel
billingPolicyToCreate := billingPolicyCreateDto{
    BillingInstrument: BillingInstrumentDto{...},
    Location: ...,
    Name: ...
}
```

## Fix

- Add inline comments describing the struct fields, or
- Add GoDoc-style comments above each model/type definition, or
- Add a `// See [X] for definition` comment referencing where the struct is defined

**Example:**

```go
// BillingPolicyResourceModel holds all fields for schema state/plan.
// See internal/services/licensing/models.go for struct definition.
var plan *BillingPolicyResourceModel
```

---

**This issue should be saved in:**
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_billing_policy_model_naming_low.md`
