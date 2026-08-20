// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

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
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/validators"
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
			"The assignment is scoped by which identifier you set. Set `environment_id` to scope it to an environment, " +
			"or `environment_group_id` to scope it to an environment group. Set neither and the assignment applies to the whole tenant.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("environment_group_id")),
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(helpers.GuidRegex), "environment_id must be a guid"),
					validators.OtherFieldRequiredWhenValueOf(
						path.Root("scope_type").Expression(),
						regexp.MustCompile("^"+scopeEnvironment+"$"), nil,
						"environment_id is required when scope_type is `environment`"),
				},
			},
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the environment group to scope the assignment to. Required when `scope_type` is `environment_group`, and not valid otherwise",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("environment_id")),
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(helpers.GuidRegex), "environment_group_id must be a guid"),
					validators.OtherFieldRequiredWhenValueOf(
						path.Root("scope_type").Expression(),
						regexp.MustCompile("^"+scopeEnvironmentGroup+"$"), nil,
						"environment_group_id is required when scope_type is `environment_group`"),
				},
			},
			"principal_id": schema.StringAttribute{
				MarkdownDescription: "The Microsoft Entra object ID of the principal to assign the role to. For a service principal this is the enterprise application object ID (`azuread_service_principal.x.object_id`), not the application (client) ID",
				Required:            true,
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
				MarkdownDescription: "The ID of the role definition to assign",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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

	state.PrincipalId = types.StringValue(found.PrincipalObjectId)
	state.PrincipalType = types.StringValue(found.PrincipalType)
	state.RoleDefinitionId = types.StringValue(found.RoleDefinitionId)
	state.Scope = types.StringValue(found.Scope)
	state.CreatedOn = types.StringValue(found.CreatedOn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *roleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every attribute requires replace, so there is no in-place update.
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

	parts := strings.Split(req.ID, "/")
	switch {
	case len(parts) == 1 && parts[0] != "":
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope_type"), scopeTenant)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[0])...)
	case len(parts) == 3 && parts[0] == "environments" && parts[1] != "" && parts[2] != "":
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope_type"), scopeEnvironment)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
	case len(parts) == 3 && parts[0] == "environmentGroups" && parts[1] != "" && parts[2] != "":
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scope_type"), scopeEnvironmentGroup)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_group_id"), parts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[2])...)
	default:
		resp.Diagnostics.AddError("Invalid import ID", importFormats)
	}
}

func findRoleAssignment(assignments []roleAssignmentDto, roleAssignmentId string) *roleAssignmentDto {
	for i := range assignments {
		if assignments[i].RoleAssignmentId == roleAssignmentId {
			return &assignments[i]
		}
	}
	return nil
}
