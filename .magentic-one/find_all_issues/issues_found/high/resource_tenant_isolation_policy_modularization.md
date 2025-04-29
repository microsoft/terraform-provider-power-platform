# Title
Tenant Isolation Policy resource file lacks proper separation of concerns and modular design.

## 
/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem
- The code combines multiple concerns in a single large file:
   - DTO conversion (`convertToDto`, `convertFromDto`) is handled within resource operations such as `Create` and `Update`, rather than being in a utility module or library.
   - API calls (`createOrUpdateTenantIsolationPolicy`) are intermixed with the higher-level resource lifecycle handling methods like `Create`, `Update`, etc., making it hard to test or reuse these API calls independently.
   - Error reporting logic is scattered and duplicated across different sections, increasing its maintenance burden.

## Impact
1. Increased cognitive load for developers working on the file.
2. Harder code maintenance due to the mixture of functional concerns.
3. Reduced testability:
   - Separate unit testing for DTO conversions or API logic becomes impractical.
   - Debugging requires navigating through lengthy methods, impacting efficiency.
4. Increased risk of bugs when modifying shared sections of logic due to lack of isolation between concerns.

Severity: `High`

## Location
File: `/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go`

## Code Issue
An example of combined concerns found in the `Create` function:
```go
// Convert the Terraform model to the API model
policyDto, diags := convertToDto(ctx, tenantInfo.TenantId, &plan)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
	return
}

// Create the policy
policy, err := r.Client.createOrUpdateTenantIsolationPolicy(ctx, tenantInfo.TenantId, *policyDto)
if err != nil {
	resp.Diagnostics.AddError(
		"Error creating tenant isolation policy",
		fmt.Sprintf("Could not create tenant isolation policy: %s", err.Error()),
	)
	return
}
```

## Fix
1. **DTO Conversion Modularization**: Move `convertToDto` and `convertFromDto` into a utility package dedicated to model transformations.
2. **API Logic Isolation**: Create a dedicated API client or helper module for tenant isolation policy operations such as `createOrUpdateTenantIsolationPolicy`.
3. **Error Handling Consolidation**: Extract common error handling functionality into helper methods or reusable functions.

Revised version of the `Create` function:
```go
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan TenantIsolationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current tenant ID
	tenantInfo, err := r.Client.TenantApi.GetTenant(ctx)
	if err != nil {
		helpers.ReportError(resp.Diagnostics, "Error retrieving tenant information", err)
		return
	}

	if tenantInfo.TenantId == "" {
		helpers.ReportError(resp.Diagnostics, "Tenant ID not found")
		return
	}

	// Convert the Terraform model to the API model (using a dedicated utility method)
	policyDto, diags := tenantIsolationPolicyUtils.ConvertToDTO(ctx, tenantInfo.TenantId, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the policy (using a dedicated API client)
	policy, err := r.Client.CreateOrUpdatePolicy(ctx, tenantInfo.TenantId, *policyDto)
	if err != nil {
		helpers.ReportError(resp.Diagnostics, "Error creating tenant isolation policy", err)
		return
	}

	// Convert the API response back to the Terraform model
	state, diags := tenantIsolationPolicyUtils.ConvertFromDTO(ctx, policy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Id = types.StringValue(tenantInfo.TenantId)
	state.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Created tenant isolation policy with ID %s", state.Id.ValueString()))
}
```

**Benefits of Fix**:
- Clear separation of concerns with modularized logic.
- Reduced cognitive load for developers when navigating the codebase.
- Enhanced testability and ease of maintenance.