// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access //nolint:revive // the underscored package name predates this file and matches every service in the repo

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// principalTypes are the principal kinds the RBAC API accepts. The API models this as an enum
// (RoleAssignmentPrincipalType) but publishes no values, so these were confirmed against it: a
// request carrying an unknown value is rejected while converting the body, whereas a known value
// with an unknown principal gets as far as PrincipalDoesNotExist. All three below reach the latter.
// Note that ServicePrincipal, the azurerm spelling, is NOT accepted; a service principal is
// `ApplicationUser` here.
var principalTypes = []string{"ApplicationUser", "Group", "User"}

var _ resource.Resource = &roleAssignmentResource{}
var _ resource.ResourceWithImportState = &roleAssignmentResource{}
var _ resource.ResourceWithValidateConfig = &roleAssignmentResource{}
var _ resource.ResourceWithModifyPlan = &roleAssignmentResource{}

func NewRoleAssignmentResource() resource.Resource {
	return &roleAssignmentResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "role_assignment",
		},
	}
}

func (r *roleAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *roleAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Power Platform administrative [role assignment](https://learn.microsoft.com/en-us/rest/api/power-platform/authorization/role-based-access-control). " +
			"Use this resource to assign Power Platform roles to service principals, users or groups. For Dataverse security roles inside an environment, use `powerplatform_security_role_assignment` instead.\n\n" +
			"~> The role based access control API is in [preview](https://learn.microsoft.com/en-us/power-platform/admin/security/role-based-access-control) and Microsoft does not recommend it for production use yet. Managing assignments requires the caller to hold the Power Platform Administrator Entra role or the Power Platform Role Based Access Control Administrator role.\n\n" +
			"The assignment is scoped by the required `scope_type`: `tenant`, `environment` with `environment_id`, or `environment_group` with `environment_group_id`. " +
			"Tenant scope is the broadest grant available, so it must be named explicitly.\n\n" +
			"The role is selected by exactly one of `role_definition_id` or `role_definition_name`; a name is resolved to its id at create time and the assignment is anchored on the id from then on.\n\n" +
			"The configuration identifies the relationship between a principal, a role and a scope, while the assignment id is computed. " +
			"If exactly one assignment for that relationship already exists it is adopted, and destroying the resource removes it. " +
			"If several duplicates exist the create fails: deduplicate them or import one first. " +
			"Expiring assignments are not adopted: this resource represents a non-expiring relationship, so configuring the same principal, role and scope alongside an expiring assignment creates a separate permanent assignment.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the role assignment",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope_type": schema.StringAttribute{
				MarkdownDescription: "Where the assignment applies. One of `tenant`, `environment` or `environment_group`. " +
					"`environment` requires `environment_id` and `environment_group` requires `environment_group_id`. " +
					"Tenant scope is the broadest grant available, so it must be asked for by name rather than by leaving the scope unset",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(scopeKinds...),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the environment to scope the assignment to. Required when `scope_type` is `environment`, and not valid otherwise",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("environment_group_id")),
				},
			},
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the environment group to scope the assignment to. Required when `scope_type` is `environment_group`, and not valid otherwise",
				Optional:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("environment_id")),
				},
			},
			"principal_id": schema.StringAttribute{
				MarkdownDescription: "The Microsoft Entra object ID of the principal to assign the role to, for every principal type. " +
					"For a service principal this is the enterprise application object ID (`azuread_service_principal.x.object_id`), not the application (client) ID. " +
					"For a user it is the user's object ID, not the email address: the API models this field as a GUID, so an email is rejected before any lookup happens",
				Required:   true,
				CustomType: customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principal_type": schema.StringAttribute{
				MarkdownDescription: "The kind of principal being assigned the role. One of `ApplicationUser` for a service principal or managed identity, `Group` for a security enabled Microsoft Entra group, or `User` for a person",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(principalTypes...),
				},
			},
			"role_definition_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the role definition to assign. Exactly one of `role_definition_id` or `role_definition_name` must be set; the id is computed when the name is used",
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.UUIDType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("role_definition_id"), path.MatchRoot("role_definition_name")),
				},
			},
			"role_definition_name": schema.StringAttribute{
				MarkdownDescription: "The name of the role definition to assign, matched case-insensitively and resolved to its id at create time. The assignment is anchored on the id from then on, so a name edit updates in place: if the new name resolves to the same role nothing changes remotely, and if it resolves to a different role the update fails rather than silently replacing the grant, since replacing by name could destroy an adopted assignment. Change `role_definition_id` to move the assignment. Replacing a name-selected assignment (for example by changing the principal) is refused with the stored id handed over, since a replacement would re-resolve the name; a taint or -replace is an explicit recreate and resolves the name afresh. Exactly one of `role_definition_id` or `role_definition_name` must be set; the name is filled in from the catalogue when the id is used, on the read after the create",
				Optional:            true,
				Computed:            true,
				CustomType:          customtypes.CaseInsensitiveStringType{},
				PlanModifiers: []planmodifier.String{
					caseInsensitiveUnchanged{},
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "The fully qualified scope of the role assignment (computed by the API)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_on": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the role assignment was created",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// ModifyPlan protects a name-selected assignment's stored role id at the planning boundary.
// Update alone cannot: a replacement (a changed principal or scope, a taint, -replace) bypasses
// Update and runs Create, and Terraform plans a replacement's new instance from the configuration
// alone, so neither state nor plan values nor private state reach that create. Re-resolving the
// unchanged name there against a drifted catalogue would silently retarget the grant. So:
//
//   - a replacement of a name-selected assignment is refused at plan time, with the stored id
//     handed over so the replacement can be applied with role_definition_id explicitly;
//   - a changed name combined with a replacement is refused, since the rename must be applied on
//     its own first;
//   - a changed name alone is validated against the stored id at plan time when the catalogue
//     answers, with Update keeping the authoritative check at apply time;
//   - a taint or -replace carries no attribute change to detect, so it re-resolves the name like
//     the fresh create it is: an explicit operator recreate of this exact resource.
func (r *roleAssignmentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// A destroy needs no plan changes.
		return
	}
	if req.State.Raw.IsNull() {
		// A fresh create, or the create leg of a replacement, which Terraform plans from the
		// configuration alone: nothing from the old instance reaches this pass, so the safety for
		// replacements lives in the prior-state pass below, which refuses to plan them for a
		// name-selected assignment in the first place.
		return
	}

	var config, plan, state roleAssignmentResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !helpers.IsKnown(config.RoleDefinitionName) {
		// Id-selected: the configured id is the explicit anchor and the only retarget mechanism.
		return
	}

	nameChanged := !strings.EqualFold(config.RoleDefinitionName.ValueString(), state.RoleDefinitionName.ValueString())
	replacing := !strings.EqualFold(config.ScopeType.ValueString(), state.ScopeType.ValueString()) ||
		!uuidEqual(plan.EnvironmentId, state.EnvironmentId) ||
		!uuidEqual(plan.EnvironmentGroupId, state.EnvironmentGroupId) ||
		!uuidEqual(plan.PrincipalId, state.PrincipalId) ||
		plan.PrincipalType.ValueString() != state.PrincipalType.ValueString()

	if nameChanged && replacing {
		resp.Diagnostics.AddError(
			"The role name cannot change in the same apply as a replacement",
			"This change replaces the assignment and also edits role_definition_name. A replacement re-creates the assignment, and resolving a new name during it could target a different role than the one this assignment holds. Apply the name change on its own first, then the replacement.",
		)
		return
	}

	if replacing {
		// Terraform plans the replacement's new instance from the configuration alone, so the
		// stored role id cannot be carried through and the unchanged name would be re-resolved
		// against today's catalogue. If the catalogue has drifted, that silently retargets the
		// grant. The replacement therefore requires the explicit id, which the error hands over.
		resp.Diagnostics.AddError(
			"Replacing a name-selected role assignment requires the explicit role id",
			fmt.Sprintf("This change replaces the assignment, and a replacement re-resolves role_definition_name against the catalogue at apply time, which can silently target a different role than the one this assignment holds. Set role_definition_id = %q for the replacement instead of the name; switch back to the name afterwards if preferred.", state.RoleDefinitionId.ValueString()),
		)
		return
	}

	if nameChanged {
		// Fail a retargeting rename at plan time when the catalogue answers; Update stays the
		// authoritative check at apply time in case the catalogue changes in between.
		if definition, err := r.Client.ResolveRoleDefinitionByName(ctx, config.RoleDefinitionName.ValueString()); err == nil {
			if !strings.EqualFold(definition.RoleDefinitionId, state.RoleDefinitionId.ValueString()) {
				resp.Diagnostics.AddError(
					"The new role definition name resolves to a different role",
					fmt.Sprintf("%q resolves to role definition %s, but this assignment holds %s. A name edit only follows a rename of the same role. To assign a different role, set role_definition_id, which replaces the assignment.", config.RoleDefinitionName.ValueString(), definition.RoleDefinitionId, state.RoleDefinitionId.ValueString()),
				)
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Deferring the rename validation of %q to apply: %s", config.RoleDefinitionName.ValueString(), err))
		}
		return
	}

	// The name is unchanged, so the stored id is the assignment's identity through any plan.
	if helpers.IsKnown(state.RoleDefinitionId) && !helpers.IsKnown(plan.RoleDefinitionId) {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("role_definition_id"), state.RoleDefinitionId)...)
	}
}

// uuidEqual compares two uuid values case-insensitively, treating null and unknown as equal only
// to themselves.
func uuidEqual(a, b customtypes.UUID) bool {
	if !helpers.IsKnown(a) || !helpers.IsKnown(b) {
		return helpers.IsKnown(a) == helpers.IsKnown(b)
	}
	return strings.EqualFold(a.ValueString(), b.ValueString())
}

// ValidateConfig enforces that scope_type and its matching id arrive together. The per-attribute
// validators cannot cover the missing-id case, because a validator attached to an absent attribute
// never runs.
func (r *roleAssignmentResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config roleAssignmentResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validateScopeSelection(config.ScopeType, config.EnvironmentId, config.EnvironmentGroupId)...)
}

func (r *roleAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T.", req.ProviderData),
		)
		return
	}
	r.Client = newRoleBasedAccessClient(providerClient.Api)
}

func (r *roleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan roleAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := plan.assignmentScope()
	tflog.Debug(ctx, fmt.Sprintf("Creating role assignment for principal %s at %s scope", plan.PrincipalId.ValueString(), scope))

	// The role is selected by exactly one of id or name. A name must resolve to exactly one
	// definition before anything is created; an id is used as given, and its display name is left
	// null for Read to fill in afterwards. No cosmetic request runs before the POST whose state
	// matters, so nothing cosmetic can consume the create's context.
	if !helpers.IsKnown(plan.RoleDefinitionId) {
		definition, err := r.Client.ResolveRoleDefinitionByName(ctx, plan.RoleDefinitionName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to resolve the role definition named %q", plan.RoleDefinitionName.ValueString()), err.Error())
			return
		}
		plan.RoleDefinitionId = customtypes.NewUUIDValue(definition.RoleDefinitionId)
	} else if !helpers.IsKnown(plan.RoleDefinitionName) {
		plan.RoleDefinitionName = customtypes.NewCaseInsensitiveStringNull()
	}

	assignment, err := r.Client.CreateRoleAssignment(ctx, scope, roleAssignmentRequestDto{
		PrincipalObjectId: plan.PrincipalId.ValueString(),
		PrincipalType:     plan.PrincipalType.ValueString(),
		RoleDefinitionId:  plan.RoleDefinitionId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to create role assignment at %s scope", scope), err.Error())
		return
	}

	plan.Id = types.StringValue(assignment.RoleAssignmentId)
	plan.Scope = types.StringValue(assignment.Scope)
	plan.CreatedOn = types.StringValue(assignment.CreatedOn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// caseInsensitiveUnchanged keeps the state value when the planned value only differs from it by
// case. Semantic equality alone does not stop the framework planning an in-place change for a
// configured case-only edit, and any planned change on the name would run the update needlessly.
type caseInsensitiveUnchanged struct{}

func (m caseInsensitiveUnchanged) Description(context.Context) string {
	return "Treats a case-only change as no change."
}

func (m caseInsensitiveUnchanged) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m caseInsensitiveUnchanged) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if strings.EqualFold(req.PlanValue.ValueString(), req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

// nullableNameValue returns a null name when the value is empty, since an Optional and
// Computed attribute must end an apply known, and "absent" is more honest than "".
func nullableNameValue(value string) customtypes.CaseInsensitiveString {
	if value == "" {
		return customtypes.NewCaseInsensitiveStringNull()
	}
	return customtypes.NewCaseInsensitiveStringValue(value)
}

func (r *roleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state roleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := state.assignmentScope()
	tflog.Debug(ctx, fmt.Sprintf("Reading role assignment %s at %s scope", state.Id.ValueString(), scope))

	assignments, err := r.Client.ListRoleAssignments(ctx, scope)
	if err != nil {
		// A gone scope (deleted environment or environment group) took its assignments with it.
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to list role assignments at %s scope", scope), err.Error())
		return
	}

	found := findRoleAssignment(assignments, state.Id.ValueString())
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.PrincipalId = customtypes.NewUUIDValue(found.PrincipalObjectId)
	state.PrincipalType = types.StringValue(found.PrincipalType)
	state.RoleDefinitionId = customtypes.NewUUIDValue(found.RoleDefinitionId)
	state.Scope = types.StringValue(found.Scope)
	state.CreatedOn = types.StringValue(found.CreatedOn)
	// A configured or previously filled name is left alone, so a recased or renamed catalogue
	// entry cannot churn state; the id anchors the assignment either way. The name is only filled
	// in when absent, which is the import path and the id-configured path.
	if !helpers.IsKnown(state.RoleDefinitionName) {
		state.RoleDefinitionName = nullableNameValue(r.Client.RoleDefinitionNameById(ctx, found.RoleDefinitionId))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles the one in-place change: an edited role_definition_name. Names must not force
// replacement, because a replacement create adopts the existing assignment and the deposed
// instance then destroys it, silently revoking the grant under create_before_destroy. Instead the
// new name is resolved: the same role id is a pure state update, and a different role id fails
// with instructions, since moving the grant is what role_definition_id is for.
func (r *roleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan, state roleAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameChanged := helpers.IsKnown(plan.RoleDefinitionName) &&
		!strings.EqualFold(plan.RoleDefinitionName.ValueString(), state.RoleDefinitionName.ValueString())
	if nameChanged {
		definition, err := r.Client.ResolveRoleDefinitionByName(ctx, plan.RoleDefinitionName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Failed to resolve the role definition named %q", plan.RoleDefinitionName.ValueString()), err.Error())
			return
		}
		if !strings.EqualFold(definition.RoleDefinitionId, state.RoleDefinitionId.ValueString()) {
			resp.Diagnostics.AddError(
				"The new role definition name resolves to a different role",
				fmt.Sprintf("%q resolves to role definition %s, but this assignment holds %s. A name edit only follows a rename of the same role. To assign a different role, set role_definition_id, which replaces the assignment.", plan.RoleDefinitionName.ValueString(), definition.RoleDefinitionId, state.RoleDefinitionId.ValueString()),
			)
			return
		}
	}

	plan.Id = state.Id
	plan.RoleDefinitionId = state.RoleDefinitionId
	plan.Scope = state.Scope
	plan.CreatedOn = state.CreatedOn
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *roleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state roleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := state.assignmentScope()
	tflog.Debug(ctx, fmt.Sprintf("Deleting role assignment %s at %s scope", state.Id.ValueString(), scope))

	if err := r.Client.DeleteRoleAssignment(ctx, scope, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to delete role assignment at %s scope", scope), err.Error())
		return
	}
}

// ImportState accepts the scope and the assignment id, mirroring the API path shape:
//
//	{role_assignment_id}                                   tenant scope
//	environments/{environment_id}/{role_assignment_id}      environment scope
//	environmentGroups/{group_id}/{role_assignment_id}        environment group scope
func (r *roleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	const importFormats = "Import ID must be one of: `{role_assignment_id}`, " +
		"`environments/{environment_id}/{role_assignment_id}`, or " +
		"`environmentGroups/{environment_group_id}/{role_assignment_id}`"

	guidRegex := regexp.MustCompile(helpers.GuidRegex)
	isGuid := func(v string) bool { return guidRegex.MatchString(v) }
	parts := strings.Split(req.ID, "/")
	switch {
	case len(parts) == 1 && isGuid(parts[0]):
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope_type"), scopeTenant)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[0])...)
	case len(parts) == 3 && parts[0] == "environments" && isGuid(parts[1]) && isGuid(parts[2]):
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope_type"), scopeEnvironment)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
	case len(parts) == 3 && parts[0] == "environmentGroups" && isGuid(parts[1]) && isGuid(parts[2]):
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope_type"), scopeEnvironmentGroup)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_group_id"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
	default:
		resp.Diagnostics.AddError("Invalid import ID", importFormats+". Every id segment must be a guid")
	}
}

// findRoleAssignment compares ids case-insensitively, since an imported id may not match the
// service's guid casing.
func findRoleAssignment(assignments []roleAssignmentDto, roleAssignmentId string) *roleAssignmentDto {
	for i := range assignments {
		if strings.EqualFold(assignments[i].RoleAssignmentId, roleAssignmentId) {
			return &assignments[i]
		}
	}
	return nil
}
