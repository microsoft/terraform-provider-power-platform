# Merged Duplicated Constants & Literals Issues

This file contains all the duplicated constants and literals issues found in the codebase, merged into a single document for easier review and management.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go`

### Title

API Version Value Hard-Coded in Each Function

### Problem

The API version value `"2020-10-01"` is duplicated as a string literal in multiple places. It should be a constant to make version management easier.

### Impact

Low to Medium. Maintainability issue: updating version means touching every occurrence, and it is easy to miss one.

### Location

Occurs in each instance of `url.Values{...}` for constructing requests.

### Code Issue

```go
// Example
RawQuery: url.Values{
 constants.API_VERSION_PARAM: []string{"2020-10-01"},
}.Encode(),
```

### Fix

Define and use a constant, probably in `constants` package:

```go
// In constants package
const ADMIN_MANAGEMENT_APP_API_VERSION = "2020-10-01"

// In usage
RawQuery: url.Values{
 constants.API_VERSION_PARAM: []string{constants.ADMIN_MANAGEMENT_APP_API_VERSION},
}.Encode(),
```

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go`

### Title

Function Length/Structure & Code Duplication

### Problem

Both `LinkEnterprisePolicy` and `UnLinkEnterprisePolicy` share significant duplicated logic (URL construction, request building, API call, waiting). This makes maintenance harder and increases risk of bugs. Separation of concerns is weak; logging, error handling, and API calls are mixed.

### Impact

Severity: **medium**. Readability and maintainability are reduced; increases technical debt and future maintenance risk.

### Location

Functions `LinkEnterprisePolicy` and `UnLinkEnterprisePolicy`.

### Code Issue

```go
// LinkEnterprisePolicy
apiUrl := &url.URL{ ... }
values := url.Values{}
values.Add("api-version", "2019-10-01")
apiUrl.RawQuery = values.Encode()
linkEnterprosePolicyDto := linkEnterprosePolicyDto{
 SystemId: systemId,
}
apiResponse, err := client.Api.Execute(...)
// UnLinkEnterprisePolicy
// (Repeated code)
```

### Fix

Refactor to extract common code (e.g., URL construction, request execution) into helper functions or methods.

```go
func buildEnterprisePolicyURL(baseURL, action, envId, envType string) string {
 return fmt.Sprintf("https://%s/providers/Microsoft.BusinessAppPlatform/environments/%s/enterprisePolicies/%s/%s?api-version=2019-10-01", 
  baseURL, envId, envType, action)
}
// Then use
linkURL := buildEnterprisePolicyURL(client.Api.GetConfig().Urls.BapiUrl, "link", environmentId, environmentType)
unlinkURL := buildEnterprisePolicyURL(client.Api.GetConfig().Urls.BapiUrl, "unlink", environmentId, environmentType)
// And so on.
```

## ISSUE 3

**File:** `/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go`

### Title

Duplicated Logic for Building API URLs With Query Strings

### Problem

The process of building API URLs with query strings (scheme, host, path, setting RawQuery) is repeated throughout the file. This code duplication can lead to subtle errors or inconsistencies if logic changes in only some places. Duplicated snippets add cognitive load and violate DRY (Don't Repeat Yourself) principles.

### Impact

Severity: Medium

Duplication makes refactoring harder, inflates the codebase, and increases testing surface area. Changes in query encoding, endpoint base paths, or global API versioning may become inconsistent across the implementation.

### Location

See similar code in almost every public method, such as `GetDataverseUserBySystemUserId`, `GetDataverseUsers`, `GetDataverseUserByAadObjectId`, `GetEnvironmentUserByAadObjectId`, and others.

### Code Issue

```go
apiUrl := &url.URL{
 Scheme: constants.HTTPS,
 Host:   environmentHost,
 Path:   "/api/data/v9.2/systemusers",
}
values := url.Values{}
values.Add("$expand", "systemuserroles_association($select=roleid,name,ismanaged,_businessunitid_value)")
apiUrl.RawQuery = values.Encode()
```

### Fix

Extract to helper(s):

```go
func buildApiUrl(scheme, host, path string, query url.Values) string {
 apiUrl := &url.URL{
  Scheme: scheme,
  Host:   host,
  Path:   path,
 }
 apiUrl.RawQuery = query.Encode()
 return apiUrl.String()
}
```

This allows future global changes and makes tests/maintenance much easier.

## ISSUE 4

**File:** `/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations.go`

### Title

Inadequate Error Handling on State Get in Read

### Problem

In the `Read` method, the code retrieves state with:

```go
var state DataSourceModel
resp.State.Get(ctx, &state)
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
 return
}
```

This calls `resp.State.Get` twice: once without checking its error and again with appending diagnostics. If the first call populates `state` or fails, the second call could yield different results or diagnostics could be duplicated/misleading. It also violates the pattern of always handling and appending errors for state operations.

### Impact

This could lead to errors being missed, diagnostics being inconsistent or duplicated, and makes the control flow harder to understand and debug. **Severity: Medium**

### Location

Method: `Read`, lines around:

```go
var state DataSourceModel
resp.State.Get(ctx, &state)
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
 return
}
```

### Code Issue

```go
var state DataSourceModel
resp.State.Get(ctx, &state)
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
 return
}
```

### Fix

Only call `Get` once, append any diagnostics, and proceed based on error state:

```go
var state DataSourceModel
diags := resp.State.Get(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
 return
}
```

## ISSUE 5

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages.go`

### Title

Missing State Reinitialization Before Accumulating Applications List

### Problem

In the Read method, `state.Applications` is potentially appended to on each Read. Since Terraform framework may preserve state between calls (depending on usage and errors), not clearing or reinitializing this list can result in duplicated or stale data if Read is called multiple times without resetting state.

### Impact

- Medium: Could cause duplicated entries in the returned list upon multiple reads.
- Data consistency issues for users.

### Location

- Read method, before iterating and appending to `state.Applications`.

### Code Issue

```go
for _, application := range applications {
 // ...
 state.Applications = append(state.Applications, app)
}
```

### Fix

**Explicitly reinitialize `state.Applications` before appending new values:**

```go
state.Applications = []TenantApplicationPackageDataSourceModel{}

for _, application := range applications {
 // ...
 state.Applications = append(state.Applications, app)
}
```

## ISSUE 6

**File:** `/workspaces/terraform-provider-power-platform/internal/services/connection/models.go`

### Title

Code Structure and Redundancy

### Problem

There is noticeable repetition in the definition of data source and resource models for shares and connections. For example, `SharesDataSourceModel` and `SharesListDataSourceModel` follow a pattern that's duplicated by `ShareResourceModel`, `SharePrincipalResourceModel`, and their connection counterparts. These models could be parameterized or further composed to reduce boilerplate and the likelihood of inconsistencies as the codebase evolves.

### Impact

Low to Medium. Code repetition increases the risk of inconsistencies, makes updates more error-prone, and adds unnecessary maintenance overhead.

### Location

Widespread repetition in model definitions for shares and connections in:

- `SharesListDataSourceModel`
- `SharesDataSourceModel`
- `ConnectionsListDataSourceModel`
- `ConnectionsDataSourceModel`
- `ShareResourceModel`
- `SharePrincipalResourceModel`
- `ResourceModel`

### Code Issue

```go
// Example of pattern repetition
type SharesListDataSourceModel struct {
 Timeouts      timeouts.Value          `tfsdk:"timeouts"`
 EnvironmentId types.String            `tfsdk:"environment_id"`
 ConnectorName types.String            `tfsdk:"connector_name"`
 ConnectionId  types.String            `tfsdk:"connection_id"`
 Shares        []SharesDataSourceModel `tfsdk:"shares"`
}

type ConnectionsListDataSourceModel struct {
 Timeouts      timeouts.Value               `tfsdk:"timeouts"`
 EnvironmentId types.String                 `tfsdk:"environment_id"`
 Connections   []ConnectionsDataSourceModel `tfsdk:"connections"`
}
```

### Fix

Investigate whether generics, embedding, or composition can help reduce redundancy. For example, extract shared fields into base structs or use type embedding:

```go
type ListDataSourceModelBase struct {
 Timeouts      timeouts.Value `tfsdk:"timeouts"`
 EnvironmentId types.String   `tfsdk:"environment_id"`
}

type SharesListDataSourceModel struct {
 ListDataSourceModelBase
 ConnectorName types.String            `tfsdk:"connector_name"`
 ConnectionId  types.String            `tfsdk:"connection_id"`
 Shares        []SharesDataSourceModel `tfsdk:"shares"`
}

type ConnectionsListDataSourceModel struct {
 ListDataSourceModelBase
 Connections []ConnectionsDataSourceModel `tfsdk:"connections"`
}
```

Adopting such refactoring will decrease maintenance effort and improve code clarity.

## ISSUE 7

**File:** `internal/services/licensing/resource_billing_policy_environment.go`

### Title

Inconsistent error handling and repeated error messages

### Problem

Throughout the CRUD methods (`Create`, `Read`, `Update`, `Delete`), the error handling uses similar but duplicated error message patterns for adding diagnostic errors to the `resp.Diagnostics` field. There are repeated blocks that capture an error, format a fairly generic error message, and return. However, there is some inconsistency, such as not logging detailed context, potentially leaking sensitive error details to the end user, and poor separation of error construction.

Additionally, in the `Read` function in particular, there's a specific case where a not-found error is handled (removes resource from state), but for all other error cases, it seemingly just passes `err.Error()` to the diagnostics. This can sometimes expose internal API or implementation error messages directly to users rather than a sanitized, high-level description.

### Impact

Severity: medium

This pattern impacts maintainability (significant repeated code), introduces potential for inconsistency in user-facing error messages, and could lead to unintentional leakage of internal errors. In the worst case, sensitive/internal details may be exposed to users if `err.Error()` is not sanitized beforehand.

### Location

Example pattern that repeats in most CRUD methods, e.g.:

```go
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
 return
}
```

### Code Issue

```go
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
 return
}
```

(and similar cases throughout the file)

### Fix

Centralize error handling and logging, use a helper function to generate error diagnostics, optionally sanitize error messages before exposing to end users. Also, ensure consistent error messages in CRUD operations.

```go
// Create a helper function for error handling
func addClientError(diags *resource.Diagnostics, action, typeName string, err error) {
 // Optionally, sanitize or wrap err.Error()
 diags.AddError(fmt.Sprintf("Client error when %s %s", action, typeName), err.Error())
}

// Usage in CRUD methods:
if err != nil {
 addClientError(&resp.Diagnostics, "updating", r.FullTypeName(), err)
 return
}
```

Additionally, consider mapping error responses to more user-friendly messages or even error codes where possible (such as for not found cases).

## ISSUE 8

**File:** `/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go`

### Title

Duplicate logic in Create and Update methods

### Problem

The Create and Update methods share a very large block of nearly-identical logic: they both fetch valid solution checker rules, validate user overrides, construct a DTO, call EnableManagedEnvironment, and then reconstruct state from the returned environment. This violates DRY (Don't Repeat Yourself) and reduces maintainability—any fix/enhancement to this logic must be duplicated, increasing the risk of subtle bugs or divergence.

### Impact

Medium. Maintenance burden increases and bug fixes can be accidentally applied in only one place, causing inconsistent behavior between resource creation and update.

### Location

Both Create and Update, from solution checker rules validation through to setting state from the API response.

### Code Issue

Example of duplicated logic:

```go
// Fetch the available solution checker rules
validRules, err := r.ManagedEnvironmentClient.FetchSolutionCheckerRules(ctx, plan.EnvironmentId.ValueString())
if err != nil {
    resp.Diagnostics.AddError("Failed to fetch solution checker rules", err.Error())
    return
}
// ...
// Validate the provided solutionCheckerRuleOverrides
// ...
// Construct DTO
// ...
// Set state from env
```

### Fix

Extract the shared logic to helper functions, e.g.:

```go
type SolutionCheckerResult struct {
    Dto environment.GovernanceConfigurationDto
    RuleOverrides *string
}

func (r *ManagedEnvironmentResource) validateAndBuildGovernanceDTO(ctx context.Context, plan *ManagedEnvironmentResourceModel) (*SolutionCheckerResult, diag.Diagnostics) {
    // move shared code here
}
```

Call the helper from both Create and Update. This ensures single-responsibility and centralization for bug fixing or future code expansion.

## ISSUE 9

**File:** `/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go`

### Title

Duplicated Null-Check Logic for PlanValue Assignment

### Problem

The assignment logic for `resp.PlanValue` to `types.StringNull()` is duplicated (i.e., both when `settingsFile.IsNull()` and `settingsFile.IsUnknown()`). This could be simplified to reduce redundancy in the code.

### Impact

While not a severe logic bug, it increases maintenance overhead and slightly reduces code clarity. Severity: **low**.

### Location

Within `PlanModifyString`:

```go
if settingsFile.IsNull() {
 resp.PlanValue = types.StringNull()
} else if settingsFile.IsUnknown() {
 resp.PlanValue = types.StringNull()
} else {
 // ...
}
```

### Code Issue

```go
if settingsFile.IsNull() {
 resp.PlanValue = types.StringNull()
} else if settingsFile.IsUnknown() {
 resp.PlanValue = types.StringNull()
} else {
 // ...
}
```

### Fix

Combine the two checks using a logical OR to reduce duplicate code:

```go
if settingsFile.IsNull() || settingsFile.IsUnknown() {
 resp.PlanValue = types.StringNull()
} else {
 // ...
}
```

---

**Total Issues Found:** 9

**Summary:**

- Low severity: 3 issues
- Medium severity: 6 issues
- High severity: 0 issues

**Primary Categories:**

- API version hardcoding and constant duplication
- Code duplication in CRUD operations
- Redundant model structures
- Error handling patterns
- State management issues
