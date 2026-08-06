// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import (
	"context"
	"fmt"

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

var _ resource.Resource = &roleBasedAccessAssignmentResource{}
var _ resource.ResourceWithImportState = &roleBasedAccessAssignmentResource{}

func NewRoleBasedAccessAssignmentResource() resource.Resource {
	return &roleBasedAccessAssignmentResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "role_based_access_assignment",
		},
	}
}

func (r *roleBasedAccessAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *roleBasedAccessAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a tenant-level [role assignment](https://learn.microsoft.com/en-us/rest/api/power-platform/authorization/role-based-access-control) in Power Platform. " +
			"Use this resource to assign roles to service principals or users at the tenant scope.",
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

func (r *roleBasedAccessAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *roleBasedAccessAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan roleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating tenant-level role assignment for principal %s", plan.EnterpriseApplicationObjectId.ValueString()))

	request := roleAssignmentRequestDto{
		PrincipalObjectId: plan.EnterpriseApplicationObjectId.ValueString(),
		PrincipalType:     plan.PrincipalType.ValueString(),
		RoleDefinitionId:  plan.RoleDefinitionId.ValueString(),
	}

	assignment, err := r.Client.CreateRoleAssignment(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create role assignment", err.Error())
		return
	}

	plan.Id = types.StringValue(assignment.RoleAssignmentId)
	plan.Scope = types.StringValue(assignment.Scope)
	plan.CreatedOn = types.StringValue(assignment.CreatedOn)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *roleBasedAccessAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state roleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading tenant-level role assignment %s", state.Id.ValueString()))

	assignments, err := r.Client.ListRoleAssignments(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list role assignments", err.Error())
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

func (r *roleBasedAccessAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes require replace — no in-place update possible
}

func (r *roleBasedAccessAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state roleBasedAccessAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleting tenant-level role assignment %s", state.Id.ValueString()))

	err := r.Client.DeleteRoleAssignment(ctx, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete role assignment", err.Error())
		return
	}
}

func (r *roleBasedAccessAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
