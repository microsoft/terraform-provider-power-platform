# Framework Utility Issues - Input Validation

This document contains all framework utility-related input validation issues found in the terraform-provider-power-platform codebase.


## ISSUE 1

# Title
Missing input validation for scopes argument throughout authentication methods

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
Many authentication methods accept `scopes []string` as input, but there is no verification that `scopes` is non-nil and contains valid, non-empty scope URIs. Passing an empty or malformed scope list may lead to hard-to-debug failures downstream or requests for access tokens without correct audience.

## Impact
Medium. If invalid scope input is accepted, may result in confusing Azure SDK/network/server errors or incorrect authentication behavior, reducing robustness for users and making debugging harder.

## Location
Throughout all auth methods including:
- AuthenticateClientCertificate
- AuthenticateUsingCli
- AuthenticateClientSecret
- AuthenticateOIDC
- AuthenticateUserManagedIdentity
- AuthenticateSystemManagedIdentity
- AuthenticateAzDOWorkloadIdentityFederation
- and indirectly through `GetTokenForScopes`

## Code Issue
No validation for argument:
```go
func (client *Auth) AuthenticateClientCertificate(ctx context.Context, scopes []string) (string, time.Time, error) {
    // ...
}
```

## Fix
Validate input at public API boundaries:

```go
if len(scopes) == 0 {
    return "", time.Time{}, errors.New("at least one scope is required for token request")
}
```
And document the behavior in GoDoc comments for each relevant method.

## ISSUE 2

# Title

Potential Redundancy and Documentation Issues in Cloud Environment Constants

##

/workspaces/terraform-provider-power-platform/internal/constants/constants.go

## Problem

The code block has several large constant groups that follow a pattern for each cloud environment (PUBLIC, USDOD, USGOV, USGOVHIGH, CHINA, EX, RX), each with a set of related endpoints. The visual pattern is useful, but there's a risk that new regions could be added incorrectly, details might diverge from actual product environments, or documentation could drift from the code, as becomes evident with the comments and the table at the top going out of sync with the actual constants.

Also, having large, repetitive constant blocks can make maintenance hard and increases risk of copy-paste errors (some are already present). There's no central data structure to validate that all required endpoints are present for any new cloud, and no comments at the constant level to clarify the mapping to "Clouds" (besides the initial table which can get out of date).

## Impact

Low severity but important for maintainability. As cloud topologies and endpoints evolve, this bulk-of-constants structure encourages gradual decay and makes errors (e.g. copy/paste) or gaps in coverage more likely. Lack of in-place documentation means new contributors may struggle to determine which constant belongs to which environment, and future additions may be inconsistent.

## Location

The constant groups for each cloud (blocks beginning with `PUBLIC_`, `USDOD_`, etc) and the undocumented mapping between the table and the constants.

## Code Issue

```go
const (
	PUBLIC_ADMIN_POWER_PLATFORM_URL     = "api.admin.powerplatform.microsoft.com"
	PUBLIC_OAUTH_AUTHORITY_URL          = "https://login.microsoftonline.com/"
	// ...
)
const (
	USDOD_ADMIN_POWER_PLATFORM_URL     = "api.admin.appsplatform.us"
	// ...
)
...
```
(Likewise for all other regions.)

## Fix

- Add explicit and machine-readable mappings between cloud codes and their endpoint sets, e.g. a struct or a map rather than repeated constant groups.
- Keep region table documentation directly above each relevant constant block, or generate documentation automatically from a data structure.
- For more maintainable code, consider something like:
  ```go
  type CloudEndpoints struct {
      AdminPowerPlatformURL     string
      OAuthAuthorityURL         string
      ...
  }

  var CloudEnvironments = map[string]CloudEndpoints{
      "Public": { ... },
      "USDoD": { ... },
      ...
  }
  ```
- Add comments on each block or field indicating how it maps to either the official region list or the published documentation.

---

**Save location:**  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/constants.go-cloud_block_structure-low.md


## ISSUE 3

# Title

Function `convertToAttrValueConnectorsGroup` prematurely returns from loop

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem

In the function `convertToAttrValueConnectorsGroup`, the function immediately returns a value upon finding the first connectors group with a matching classification. If there are multiple connector groups with the same classification, only the first one is included and the rest are ignored. This could lead to missed data if the input slice contains more than one group of a given classification.

## Impact

Medium severity. This might cause incomplete data to be returned if more than one connectors group with the same classification is present, and thus can lead to data loss or unexpected behavior.

## Location

Lines 80-87:

```go
func convertToAttrValueConnectorsGroup(classification string, connectorsGroup []dlpConnectorGroupsModelDto) basetypes.SetValue {
	var connectorValues []attr.Value
	for _, conn := range connectorsGroup {
		if conn.Classification == classification {
			return types.SetValueMust(connectorSetObjectType, convertToAttrValueConnectors(conn, connectorValues))
		}
	}
	return types.SetValueMust(connectorSetObjectType, []attr.Value{})
}
```

## Code Issue

```go
for _, conn := range connectorsGroup {
	if conn.Classification == classification {
		return types.SetValueMust(connectorSetObjectType, convertToAttrValueConnectors(conn, connectorValues))
	}
}
```

## Fix

Accumulate all matching connector groups and return them together. For example:

```go
func convertToAttrValueConnectorsGroup(classification string, connectorsGroup []dlpConnectorGroupsModelDto) basetypes.SetValue {
	var connectorValues []attr.Value
	for _, conn := range connectorsGroup {
		if conn.Classification == classification {
			connectorValues = append(connectorValues, convertToAttrValueConnectors(conn, []attr.Value{})...)
		}
	}
	return types.SetValueMust(connectorSetObjectType, connectorValues)
}
```


## ISSUE 4

# Issue 1: Potential Control Flow and Nil Pointer Handling in Query Appending

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go

## Problem

The `appendQuery` function appends non-nil query parts to the OData query, but there is no handling for what would happen if the input pointer (`query`) itself is nil. Additionally, as this function directly operates on the string pointer, improper/misuse or accidentally providing a nil value could lead to panics. 

## Impact

If a nil pointer is passed as the base query to `appendQuery`, a runtime panic will be triggered (`invalid memory address or nil pointer dereference`). Severity is **high** as this is a potential runtime crash.

## Location

Line(s) 81-92

## Code Issue

```go
func appendQuery(query, part *string) {
	if part != nil {
		if len(*query) > 0 {
			*query += "&"
		}
		*query = strings.Join([]string{*query, *part}, "")
	}
}
```

## Fix

Add nil check for `query` pointer and consider returning an error or avoiding mutation if `query` is nil.

```go
func appendQuery(query, part *string) {
	if query == nil {
		// avoid panic and/or log error
		return
	}
	if part != nil {
		if len(*query) > 0 {
			*query += "&"
		}
		*query = strings.Join([]string{*query, *part}, "")
	}
}
```


## ISSUE 5

# No validation or documentation on exported struct fields

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

The struct `OtherFieldRequiredWhenValueOfValidator` is exported and its fields are also exported, but there is no documentation (`godoc`) for its fields. Moreover, there is no built-in validation/sanitization on values assigned to these fields, which might lead to improper initialization or misuse.

## Impact

Lack of documentation reduces usability by other developers and can lead to out-of-range or invalid values being set, resulting in unclear code behavior or latent bugs. This is particularly important for code forming a reusable API component. Severity: **low**.

## Location

```go
type OtherFieldRequiredWhenValueOfValidator struct {
	OtherFieldExpression   path.Expression
	OtherFieldValueRegex   *regexp.Regexp
	CurrentFieldValueRegex *regexp.Regexp
	ErrorMessage           string
}
```

## Code Issue

```go
type OtherFieldRequiredWhenValueOfValidator struct {
	OtherFieldExpression   path.Expression
	OtherFieldValueRegex   *regexp.Regexp
	CurrentFieldValueRegex *regexp.Regexp
	ErrorMessage           string
}
```

## Fix

- Add Go doc comments to each exported field describing their usage.
- Consider either making fields unexported (if not for public usage) or providing a constructor that validates input.

Example:

```go
// OtherFieldRequiredWhenValueOfValidator validates that another field is present or matches a value when a certain condition is true.
type OtherFieldRequiredWhenValueOfValidator struct {
	// OtherFieldExpression is the path expression to the other required field.
	OtherFieldExpression path.Expression

	// OtherFieldValueRegex is the regex to match the other field's value.
	OtherFieldValueRegex *regexp.Regexp

	// CurrentFieldValueRegex is the regex to match the current field's value.
	CurrentFieldValueRegex *regexp.Regexp

	// ErrorMessage is the message shown when validation fails.
	ErrorMessage string
}
```

---

This issue impacts code structure, readability and maintainability, and should be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/other_field_required_when_value_of_validator_structure_low_3.md`


## ISSUE 6

# Lack of Validation for Struct Inputs

##

/workspaces/terraform-provider-power-platform/internal/helpers/typeinfo.go

## Problem

The `TypeInfo` struct does not provide any input validation when creating instances. For example, it is possible to create an object with an empty `TypeName`, which would lead to an invalid type string like `powerplatform_` in `FullTypeName`. Having such invalid names could propagate hidden errors in larger contexts.

## Impact

Medium. Absence of validation could result in invalid resource or data source type names, possibly causing issues during downstream operations or harming user experience.

## Location

```go
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}
```
...
```go
func (t *TypeInfo) FullTypeName() string {
	if t.ProviderTypeName == "" {
		return fmt.Sprintf("powerplatform_%s", t.TypeName)
	}

	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName)
}
```

## Code Issue

```go
func (t *TypeInfo) FullTypeName() string {
	if t.ProviderTypeName == "" {
		return fmt.Sprintf("powerplatform_%s", t.TypeName)
	}

	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName)
}
```

## Fix

Add validation for `TypeName` when constructing `TypeInfo` or running `FullTypeName`, and return an error if it’s missing or invalid.

```go
func (t *TypeInfo) FullTypeName() (string, error) {
	if t.TypeName == "" {
		return "", fmt.Errorf("TypeName must not be empty")
	}

	if t.ProviderTypeName == "" {
		return fmt.Sprintf("powerplatform_%s", t.TypeName), nil
	}

	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName), nil
}
```
Or enforce TypeInfo creation only via constructor.


## ISSUE 7

# Potential for Index Out of Range Panic Due to Inadequate Input Validation

##

/workspaces/terraform-provider-power-platform/internal/helpers/uri.go

## Problem

The functions `BuildEnvironmentHostUri` and `BuildTenantHostUri` lack input validation for the incoming `environmentId` or `tenantId` strings. These functions assume that, after removing hyphens, the identifier is at least two characters long. If a caller inadvertently or maliciously passes a shorter string, indexing with `envId[len(envId)-2:]` and `envId[:len(envId)-2]` will cause a runtime panic due to an "index out of range" error.

## Impact

**Severity: High**

A panic can terminate the overall process (such as a provider or an automation pipeline) unexpectedly, leading to unreliable software behavior. This is particularly critical in libraries or modules consumed by external layers, such as Terraform providers, where input may originate from user configuration.

## Location

Lines inside:
- `BuildEnvironmentHostUri`
- `BuildTenantHostUri`

## Code Issue

```go
envId := strings.ReplaceAll(environmentId, "-", "")
realm := string(envId[len(envId)-2:])
envId = envId[:len(envId)-2]
```

and

```go
envId := strings.ReplaceAll(tenantId, "-", "")
realm := string(envId[len(envId)-2:])
envId = envId[:len(envId)-2]
```

## Fix

Add input validation to ensure the processed ID has at least two characters, returning an empty string or error if not.

```go
func BuildEnvironmentHostUri(environmentId, powerPlatformUrl string) string {
	envId := strings.ReplaceAll(environmentId, "-", "")
	if len(envId) < 2 {
		// Optionally, log or handle the error accordingly.
		return ""
	}
	realm := envId[len(envId)-2:]
	envId = envId[:len(envId)-2]
	return fmt.Sprintf("%s.%s.environment.%s", envId, realm, powerPlatformUrl)
}

func BuildTenantHostUri(tenantId, powerPlatformUrl string) string {
	envId := strings.ReplaceAll(tenantId, "-", "")
	if len(envId) < 2 {
		// Optionally, log or handle the error accordingly.
		return ""
	}
	realm := envId[len(envId)-2:]
	envId = envId[:len(envId)-2]
	return fmt.Sprintf("%s.%s.tenant.%s", envId, realm, powerPlatformUrl)
}
```

---

This prevents panics by ensuring the input is sufficiently long before accessing string slices. Further handling (such as returning an error) could be used depending on the design requirements.


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
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
