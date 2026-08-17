// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &environmentGroupRoleBasedAccessAssignmentResource{}
var _ resource.ResourceWithImportState = &environmentGroupRoleBasedAccessAssignmentResource{}

func NewEnvironmentGroupRoleBasedAccessAssignmentResource() resource.Resource {
	return &environmentGroupRoleBasedAccessAssignmentResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_group_role_based_access_assignment",
		},
	}
}

func (r *environmentGroupRoleBasedAccessAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *environmentGroupRoleBasedAccessAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a [role assignment](https://learn.microsoft.com/en-us/rest/api/power-platform/authorization/role-based-access-control) scoped to an environment group in Power Platform. " +
			"Use this resource to assign roles to service principals or users at the environment group level.",
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
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the environment group",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enterprise_application_object_id": schema.StringAttribute{
				MarkdownDescription: "The object ID of the enterprise application (service principal) or user to assign the role to. For `ApplicationUser` principals this is the enterprise application object ID, not the application (client) ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principal_type": schema.StringAttribute{
				MarkdownDescription: "The type of principal (e.g., `ApplicationUser`, `User`)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
				MarkdownDescription: "The scope of the role assignment (computed by the API)",
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

func (r *environmentGroupRoleBasedAccessAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *environmentGroupRoleBasedAccessAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan environmentGroupRoleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envGroupId := plan.EnvironmentGroupId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Creating role assignment for principal %s in environment group %s", plan.EnterpriseApplicationObjectId.ValueString(), envGroupId))

	request := roleAssignmentRequestDto{
		PrincipalObjectId: plan.EnterpriseApplicationObjectId.ValueString(),
		PrincipalType:     plan.PrincipalType.ValueString(),
		RoleDefinitionId:  plan.RoleDefinitionId.ValueString(),
	}

	assignment, err := r.Client.CreateEnvironmentGroupRoleAssignment(ctx, envGroupId, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create environment group role assignment", err.Error())
		return
	}

	plan.Id = types.StringValue(assignment.RoleAssignmentId)
	plan.Scope = types.StringValue(assignment.Scope)
	plan.CreatedOn = types.StringValue(assignment.CreatedOn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentGroupRoleBasedAccessAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state environmentGroupRoleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envGroupId := state.EnvironmentGroupId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading role assignment %s in environment group %s", state.Id.ValueString(), envGroupId))

	assignments, err := r.Client.ListEnvironmentGroupRoleAssignments(ctx, envGroupId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list environment group role assignments", err.Error())
		return
	}

	var found *roleAssignmentDto
	for i := range assignments {
		if assignments[i].RoleAssignmentId == state.Id.ValueString() {
			found = &assignments[i]
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.EnterpriseApplicationObjectId = types.StringValue(found.PrincipalObjectId)
	state.PrincipalType = types.StringValue(found.PrincipalType)
	state.RoleDefinitionId = types.StringValue(found.RoleDefinitionId)
	state.Scope = types.StringValue(found.Scope)
	state.CreatedOn = types.StringValue(found.CreatedOn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentGroupRoleBasedAccessAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes require replace — no in-place update possible
}

func (r *environmentGroupRoleBasedAccessAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state environmentGroupRoleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envGroupId := state.EnvironmentGroupId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting role assignment %s in environment group %s", state.Id.ValueString(), envGroupId))

	err := r.Client.DeleteEnvironmentGroupRoleAssignment(ctx, envGroupId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete environment group role assignment", err.Error())
		return
	}
}

func (r *environmentGroupRoleBasedAccessAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: "environment_group_id/role_assignment_id"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in the format: environment_group_id/role_assignment_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_group_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
