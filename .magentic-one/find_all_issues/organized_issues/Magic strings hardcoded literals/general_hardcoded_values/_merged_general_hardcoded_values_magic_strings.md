# Magic Strings and Hardcoded Literals - General Hardcoded Values Issues

This document consolidates all magic strings and hardcoded literals issues found in general code files.


## ISSUE 1

# Title
Error message details are inconsistent and could be improved for OIDC credential instantiation

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The error checks in `NewOidcCredential` method return hardcoded, low-information errors (e.g., `"tenant is required for OIDC credential"`, `"request Token is required for OIDC credential"`). These are not wrapped or annotated, and do not use constants, leading to less consistency and not enabling callers to distinguish error types. Furthermore, there are places where more contextual information can be given, or a specific error type (possibly a typed error) could be beneficial to facilitate control flow for the caller.

## Impact
Medium severity. Poor error context and inconsistent error reporting can make troubleshooting difficult for users and maintainers, and hinder programmatic error handling downstream.

## Location
Within the implementation of `NewOidcCredential`:

## Code Issue
```go
if c.requestToken == "" {
    return nil, errors.New("request Token is required for OIDC credential")
}
if c.requestUrl == "" {
    return nil, errors.New("request URL is required for OIDC credential")
}
if options.TenantID == "" {
    return nil, errors.New("tenant is required for OIDC credential")
}
if options.ClientID == "" {
    return nil, errors.New("client is required for OIDC credential")
}
```

## Fix
Provide more actionable and detailed error messages, use error wrapping where applicable, and consider custom error types if errors need to be programmatically distinguished. Example:

```go
if c.requestToken == "" {
	return nil, fmt.Errorf("NewOIDCCredential: missing required parameter: requestToken")
}
if c.requestURL == "" {
	return nil, fmt.Errorf("NewOIDCCredential: missing required parameter: requestURL")
}
if options.TenantID == "" {
	return nil, fmt.Errorf("NewOIDCCredential: missing required parameter: TenantID")
}
if options.ClientID == "" {
	return nil, fmt.Errorf("NewOIDCCredential: missing required parameter: ClientID")
}
```
Or, if strict error handling is needed for the API, define specific error types.

---

## ISSUE 2

# Magic Strings and Constants

## 
/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/dto.go

## Problem

There are numerous raw, repeated string literals, especially for keys (e.g., `"ai_generative_settings"`, `"backup_retention"`, `"sharing_controls"`, etc.) and rule identifiers scattered throughout the code. These are "magic strings" and present maintainability and typo risks. Using consts/enums would make refactoring and verification easier, and avoids mistakes due to typos or changes.

## Impact

- **Severity:** Medium
- Risk of silent bugs from typos.
- Harder code refactoring, auditing, and documentation.
- Poor discoverability and easier to miss updates across the codebase.

## Location

Widespread, for example in:

```go
aiGenerativeSettingsObj := attrs["ai_generative_settings"]
...
backupRetentionObj := attrs["backup_retention"]
...
solutionCheckerObj := attrs["solution_checker_enforcement"]
...
makerWelcomeContentObj := attrs["maker_welcome_content"]
...
rule := environmentGroupRuleSetParameterDto{
    ...
    Type:             AI_GENERATED_DESC,
    ...
}
```

## Code Issue

```go
aiGenerativeSettingsObj := attrs["ai_generative_settings"]
```

## Fix

Define constants at the top of the file (or a shared package) for all such keys and identifiers:

```go
const (
    AttrAiGenerativeSettings = "ai_generative_settings"
    AttrBackupRetention      = "backup_retention"
    AttrSharingControls      = "sharing_controls"
    // ... etc.

    TypeAiGeneratedDesc      = "AI_GENERATED_DESC"
    TypeBackupRetention      = "BACKUP_RETENTION"
    // ... etc.
)
...
aiGenerativeSettingsObj := attrs[AttrAiGenerativeSettings]
...
```

This improves IDE support and reduces the chance of copy-paste/typo bugs as well as facilitates any renaming.

---

---

## ISSUE 3

# Title

Hardcoded string values for connector type

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem

The connector `Type` is hardcoded as `"Microsoft.PowerApps/apis"` in multiple places in the code rather than being extracted from a constant or from configuration. Magic strings in code can lead to spelling inconsistencies and make it harder to update or maintain the namespace/type in the future.

## Impact

Low severity. This is a maintainability and reliability problem—if the connector type ever needs to change, every occurrence has to be manually updated, increasing the risk of mistakes.

## Location

There are multiple locations, for example:

Lines 112–121 (`getConnectorGroup`):

```go
for _, connector := range connectors {
	connectorGroup.Connectors = append(connectorGroup.Connectors, dlpConnectorModelDto{
		Id:                        connector.Id.ValueString(),
		Type:                      "Microsoft.PowerApps/apis", // <--- hardcoded
		DefaultActionRuleBehavior: connector.DefaultActionRuleBehavior.ValueString(),
		ActionRules:               convertToDlpActionRule(connector),
		EndpointRules:             convertToDlpEndpointRule(connector),
	})
}
```

Lines 137 and 147 (`convertToDlpConnectorGroup`):

```go
connectorGroup.Connectors = append(connectorGroup.Connectors, dlpConnectorModelDto{
	Id:   connector.Id.ValueString(),
	Type: "Microsoft.PowerApps/apis", // <--- hardcoded
	DefaultActionRuleBehavior: defaultAction,
	ActionRules:               convertToDlpActionRule(connector),
	EndpointRules:             convertToDlpEndpointRule(connector),
})
```

## Code Issue

```go
Type: "Microsoft.PowerApps/apis",
```

## Fix

Extract the type string into a package-level constant, and use that constant throughout the codebase. For example:

```go
const connectorTypePowerApps = "Microsoft.PowerApps/apis"
```

And then in all code locations:

```go
Type: connectorTypePowerApps,
```

---

## ISSUE 4

# Issue 4: Repeated Hardcoding of String Literal

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem

The `NetworkIsolation` field is hardcoded as `"PublicNetwork"` in the DTO creation. If `"PublicNetwork"` is a constant value used throughout the package, it would be better to declare it as a `const` for reusability and to avoid typos.

## Impact

Severity: **Low**

Hardcoded literals lead to technical debt and increased risk of typos, especially when used in multiple locations or subject to change.

## Location

```go
NetworkIsolation:            "PublicNetwork",
```

## Fix

Define a package-level constant for the value.

```go
const DefaultNetworkIsolation = "PublicNetwork"

// ...
NetworkIsolation: DefaultNetworkIsolation,
```

---

## ISSUE 5

# Simplification and Validation: Use of Empty String for Null UUIDs

##

/workspaces/terraform-provider-power-platform/internal/services/environment/models.go

## Problem

For nullable UUID values like `BillingPolicyId` and `EnvironmentGroupId`, the code assigns `constants.ZERO_UUID` (likely a zero or empty UUID string) to represent a "null" or unset state. This is both verbose and error-prone, as it relies on magic string checks and manual handling.

## Impact

- **Severity:** Medium
- Makes code harder to follow and maintain.
- Potential for bugs if another part of the code misinterprets or mismatches the zero UUID.
- Inconsistent with Go's best practices for handling optional/nullable values.

## Location

```go
if !environmentSource.BillingPolicyId.IsNull() && environmentSource.BillingPolicyId.ValueString() != constants.ZERO_UUID {
	environmentDto.Properties.BillingPolicy = BillingPolicyDto{
		Id: environmentSource.BillingPolicyId.ValueString(),
	}
}
...
if !environmentSource.EnvironmentGroupId.IsNull() && !environmentSource.EnvironmentGroupId.IsUnknown() {
	environmentDto.Properties.ParentEnvironmentGroup = &ParentEnvironmentGroupDto{Id: environmentSource.EnvironmentGroupId.ValueString()}
}
...
func convertEnvironmentGroupFromDto(environmentDto EnvironmentDto, model *SourceModel) {
	if environmentDto.Properties.ParentEnvironmentGroup != nil {
		model.EnvironmentGroupId = types.StringValue(environmentDto.Properties.ParentEnvironmentGroup.Id)
	} else {
		model.EnvironmentGroupId = types.StringValue(constants.ZERO_UUID)
	}
}
func convertBillingPolicyModelFromDto(environmentDto EnvironmentDto, model *SourceModel) {
	if environmentDto.Properties.BillingPolicy != nil && environmentDto.Properties.BillingPolicy.Id != "" {
		model.BillingPolicyId = types.StringValue(environmentDto.Properties.BillingPolicy.Id)
	} else {
		model.BillingPolicyId = types.StringValue(constants.ZERO_UUID)
	}
}
```

## Code Issue

```go
model.BillingPolicyId = types.StringValue(constants.ZERO_UUID)
...
if !environmentSource.BillingPolicyId.IsNull() && environmentSource.BillingPolicyId.ValueString() != constants.ZERO_UUID
```

## Fix

Use `types.StringNull()` for unset/optional values, which is idiomatic in Terraform plugin development and with the `types` API.

```go
func convertEnvironmentGroupFromDto(environmentDto EnvironmentDto, model *SourceModel) {
	if environmentDto.Properties.ParentEnvironmentGroup != nil {
		model.EnvironmentGroupId = types.StringValue(environmentDto.Properties.ParentEnvironmentGroup.Id)
	} else {
		model.EnvironmentGroupId = types.StringNull()
	}
}
func convertBillingPolicyModelFromDto(environmentDto EnvironmentDto, model *SourceModel) {
	if environmentDto.Properties.BillingPolicy != nil && environmentDto.Properties.BillingPolicy.Id != "" {
		model.BillingPolicyId = types.StringValue(environmentDto.Properties.BillingPolicy.Id)
	} else {
		model.BillingPolicyId = types.StringNull()
	}
}
```

When ingesting values, check for `.IsNull()` instead of comparing to a zero/empty UUID string.

---

## ISSUE 6

# Title

Potential Inconsistent API Parameter Mapping for Connector Groups

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

The mapping between UI/group labels (“Confidential”, “General”, “Blocked”) and how those map into API-side expectations (`business_connectors`, `non_business_connectors`, etc.) is spread via magic strings throughout the codebase, especially in conversions such as `convertToAttrValueConnectorsGroup` and `convertToDlpConnectorGroup`. Mixing these display group names and schema field references can easily lead to typos and hard-to-update code, as well as mapping mismatches if API requirements or schema names change. The lack of a central mapping function or structure reduces maintainability as the product evolves.

## Impact

Severity: Medium

Medium maintainability and consistency impact. Mistyped or parity-breaking hardcoded values may introduce bugs that are hard to follow, as field semantics are not defined centrally.

## Location

```go
state.BusinessGeneralConnectors = convertToAttrValueConnectorsGroup("Confidential", policy.ConnectorGroups)
state.NonBusinessConfidentialConnectors = convertToAttrValueConnectorsGroup("General", policy.ConnectorGroups)
state.BlockedConnectors = convertToAttrValueConnectorsGroup("Blocked", policy.ConnectorGroups)
...
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.BusinessGeneralConnectors))
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.NonBusinessConfidentialConnectors))
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors))
```

## Code Issue

```go
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Confidential", plan.BusinessGeneralConnectors))
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "General", plan.NonBusinessConfidentialConnectors))
policyToCreate.ConnectorGroups = append(policyToCreate.ConnectorGroups, convertToDlpConnectorGroup(ctx, resp.Diagnostics, "Blocked", plan.BlockedConnectors))
```

## Fix

- Define a central mapping function or constant map for connector group semantics.
- Call central conversion/mapping helpers throughout resource methods.

```go
const (
	BusinessGroup     = \"Confidential\"
	NonBusinessGroup  = \"General\"
	BlockedGroup      = \"Blocked\"
)

var groupFieldMap = map[string]string{
	BusinessGroup:    \"business_connectors\",
	NonBusinessGroup: \"non_business_connectors\",
	BlockedGroup:     \"blocked_connectors\",
}
```

Then use these constants and map lookups in your conversion, resource, and validation logic.

---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_dlp_policy_magic_strings_medium.md

---

## ISSUE 7

# Title

Magic strings for settings values reduce code clarity and future expansion

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go

## Problem

In several places, hardcoded string literals like "true", "false", and settings values such as "Standard" populate the configuration DTOs and resource state. Using magic strings can reduce code clarity, increase the risk of typos, and make maintenance or refactoring harder (if the meaning or valid values change). Ideally, these string values should be constants or enums for easy replacement, reuse, and discoverability.

## Impact

Low. This impacts maintainability and code clarity, but does not cause immediate runtime errors if all strings remain valid and consistent.

## Location

Within Create, Update, and other DTO setup locations:

## Code Issue

```go
ProtectionLevel: "Standard",
IncludeOnHomepageInsights: "false",
DisableAiGeneratedDescriptions: "false",
// many similar instances ...
```

## Fix

Define constants at the top of the file:

```go
const (
    ProtectionLevelStandard = "Standard"
    IncludeOnHomepageInsightsFalse = "false"
    DisableAiGeneratedDescriptionsFalse = "false"
)
```
And reference them throughout the code. Optionally, group them under type aliases or enums for even better clarity and refactorability.

---


# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for "copilot-commit-message-instructions.md" how to write description.
- `<issue_number>` pick the issue number or PR number
```
