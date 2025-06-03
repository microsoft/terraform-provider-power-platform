# Issue: Hardcoded Strings for Policy Types (Magic Strings)

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go

## Problem

The string values `"NETWORK_INJECTION_POLICY_TYPE"` and `"ENCRYPTION_POLICY_TYPE"` are referenced multiple times directly. However, the actual values for these constants are not shown in this file, and using magic strings or undefined references can lead to maintenance challenges and risk of typos.

Additionally, the use of format specifiers in `MarkdownDescription` further obfuscates what those values are at a glance:

```go
MarkdownDescription: fmt.Sprintf("Policy type [%s, %s]", NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
```

## Impact

- Maintenance burden and risk of silent bugs if magic strings are mistyped or duplicated.
- Degrades code readability and makes it harder to update the policy types in the future.
- **Severity:** Low

## Location

```go
MarkdownDescription: fmt.Sprintf("Policy type [%s, %s]", NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
Validators: []validator.String{
    stringvalidator.OneOf(NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
},
...
if state.PolicyType.ValueString() == NETWORK_INJECTION_POLICY_TYPE ...
```

## Code Issue

```go
// Use of magic strings/constants throughout
MarkdownDescription: fmt.Sprintf("Policy type [%s, %s]", NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
Validators: []validator.String{
    stringvalidator.OneOf(NETWORK_INJECTION_POLICY_TYPE, ENCRYPTION_POLICY_TYPE),
},
if state.PolicyType.ValueString() == NETWORK_INJECTION_POLICY_TYPE ...
```

## Fix

Define and use clearly named `const` values at the top of this file:

```go
const (
    NetworkInjectionPolicyType = "NetworkInjection"
    EncryptionPolicyType = "Encryption"
)
```

Then use them consistently:

```go
MarkdownDescription: fmt.Sprintf("Policy type [%s, %s]", NetworkInjectionPolicyType, EncryptionPolicyType),
Validators: []validator.String{
    stringvalidator.OneOf(NetworkInjectionPolicyType, EncryptionPolicyType),
},
if state.PolicyType.ValueString() == NetworkInjectionPolicyType ...
```

---

This markdown will be saved under:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_enterprise_policy_structure_low_policy_constants.md`
