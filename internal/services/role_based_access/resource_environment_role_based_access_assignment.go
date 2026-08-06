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

var _ resource.Resource = &environmentRoleBasedAccessAssignmentResource{}
var _ resource.ResourceWithImportState = &environmentRoleBasedAccessAssignmentResource{}

func NewEnvironmentRoleBasedAccessAssignmentResource() resource.Resource {
	return &environmentRoleBasedAccessAssignmentResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_role_based_access_assignment",
		},
	}
}

func (r *environmentRoleBasedAccessAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *environmentRoleBasedAccessAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a [role assignment](https://learn.microsoft.com/en-us/rest/api/power-platform/authorization/role-based-access-control) scoped to an environment in Power Platform. " +
			"Use this resource to assign roles to service principals or users at the environment level.",
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
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the environment",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principal_object_id": schema.StringAttribute{
				MarkdownDescription: "The object ID of the principal (service principal or user) to assign the role to",
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

func (r *environmentRoleBasedAccessAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.Client = NewRoleBasedAccessClient(providerClient.Api)
}

func (r *environmentRoleBasedAccessAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan environmentRoleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envId := plan.EnvironmentId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Creating role assignment for principal %s in environment %s", plan.PrincipalObjectId.ValueString(), envId))

	request := roleAssignmentRequestDto{
		PrincipalObjectId: plan.PrincipalObjectId.ValueString(),
		PrincipalType:     plan.PrincipalType.ValueString(),
		RoleDefinitionId:  plan.RoleDefinitionId.ValueString(),
	}

	assignment, err := r.Client.CreateEnvironmentRoleAssignment(ctx, envId, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create environment role assignment", err.Error())
		return
	}

	plan.Id = types.StringValue(assignment.RoleAssignmentId)
	plan.Scope = types.StringValue(assignment.Scope)
	plan.CreatedOn = types.StringValue(assignment.CreatedOn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentRoleBasedAccessAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state environmentRoleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envId := state.EnvironmentId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading role assignment %s in environment %s", state.Id.ValueString(), envId))

	assignments, err := r.Client.ListEnvironmentRoleAssignments(ctx, envId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list environment role assignments", err.Error())
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

	state.PrincipalObjectId = types.StringValue(found.PrincipalObjectId)
	state.PrincipalType = types.StringValue(found.PrincipalType)
	state.RoleDefinitionId = types.StringValue(found.RoleDefinitionId)
	state.Scope = types.StringValue(found.Scope)
	state.CreatedOn = types.StringValue(found.CreatedOn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentRoleBasedAccessAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes require replace — no in-place update possible
}

func (r *environmentRoleBasedAccessAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state environmentRoleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envId := state.EnvironmentId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting role assignment %s in environment %s", state.Id.ValueString(), envId))

	err := r.Client.DeleteEnvironmentRoleAssignment(ctx, envId, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete environment role assignment", err.Error())
		return
	}
}

func (r *environmentRoleBasedAccessAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: "environment_id/role_assignment_id"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in the format: environment_id/role_assignment_id")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
