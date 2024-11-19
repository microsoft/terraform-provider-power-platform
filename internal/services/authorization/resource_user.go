// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authorization

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers/array"
)

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
	return &UserResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "user",
		},
	}
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "This resource associates a user to a Power Platform environment.\n\n" +
			"Additional Resources:\n\n" +
			"* [Add users to an environment](https://learn.microsoft.com/power-platform/admin/add-users-to-environment)\n\n" +
			"* [Overview of User Security](https://learn.microsoft.com/power-platform/admin/grant-users-access)\n\n" +
			"*Note:* When starting with non Dataverse environments, and adding Dataverse later, the 'Environment Admin' and 'Environment Maker' used earlier in `security_roles` will not work inside Dataverse. You will need to use the Dataverse security roles instead.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique user id (guid)",
				Description:         "Unique user id (guid)",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Unique environment id (guid)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"business_unit_id": schema.StringAttribute{
				MarkdownDescription: "Id of the business unit to which the user belongs",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"aad_id": schema.StringAttribute{
				MarkdownDescription: "Entra user object id",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"security_roles": schema.SetAttribute{
				MarkdownDescription: "Security roles Ids assigned to the Dataverse user" +
					"When working with non Dataverse environments, only 'Environment Admin' and 'Environment Maker' role values are allowed",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"user_principal_name": schema.StringAttribute{
				MarkdownDescription: "User principal name",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "User first name",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "User last name",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_delete": schema.BoolAttribute{
				MarkdownDescription: "Disable delete. When set to `True` is expects that (Disable Delte)[https://learn.microsoft.com/power-platform/admin/delete-users?WT.mc_id=ppac_inproduct_settings#soft-delete-users-in-power-platform] feature to be enabled." +
					"Removing resource will try to delete the systemuser from Dataverse. This is the default behaviour. If you just want to remove the resource and not delete the user from Dataverse, set this propertyto `False`\n\n" +
					"**This attribute applies only when working with dataverse users.**",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.UserClient = newUserClient(clientApi)
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasEnvDataverse, err := r.UserClient.EnvironmentHasDataverse(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Dataverse exist in eviroment %t", hasEnvDataverse))

	newUser := userDto{}
	if hasEnvDataverse {
		user, err := r.UserClient.CreateDataverseUser(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}

		user, err = r.UserClient.AddDataverseSecurityRoles(ctx, plan.EnvironmentId.ValueString(), user.Id, plan.SecurityRoles)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
		newUser = *user
	} else {
		// todo disalbe delete should be set to false.
		err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		}

		user, err := r.UserClient.CreateEnvironmentUser(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString(), plan.SecurityRoles)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}

		rolesBytes, err := json.Marshal(user.SecurityRoles)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error marshalling security roles", fmt.Sprintf("Error marshalling security roles: %s", err.Error()),
			)
			return
		}
		resp.Private.SetKey(ctx, "role", rolesBytes)

		newUser = *user
	}

	model := convertDataverseFromUserDto(&newUser, plan.DisableDelete.ValueBool())
	plan.Id = model.Id
	plan.AadId = model.AadId
	req.Plan.SetAttribute(ctx, path.Root("security_roles"), model.SecurityRoles)
	plan.UserPrincipalName = model.UserPrincipalName
	plan.FirstName = model.FirstName
	plan.LastName = model.LastName
	plan.DisableDelete = model.DisableDelete
	plan.BusinessUnitId = model.BusinessUnitId

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var state *UserResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hasEnvDataverse, err := r.UserClient.EnvironmentHasDataverse(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Dataverse exist in eviroment %t", hasEnvDataverse))

	updateUser := userDto{}
	if hasEnvDataverse {
		user, err := r.UserClient.GetDataverseUserByAadObjectId(ctx, state.EnvironmentId.ValueString(), state.AadId.ValueString())
		if err != nil {
			if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
		updateUser = *user
	} else {
		user, err := r.UserClient.GetEnvironmentUserByAadObjectId(ctx, state.EnvironmentId.ValueString(), state.AadId.ValueString())
		// if all the security roles are removed, the user will not be found
		if user.AadObjectId == "" {
			user.AadObjectId = state.AadId.ValueString()
			user.DomainName = state.UserPrincipalName.ValueString()
		}

		if err != nil {
			if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
				resp.State.RemoveResource(ctx)
				return
			}
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}

		rolesBytes, err := json.Marshal(user.SecurityRoles)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error marshalling security roles", fmt.Sprintf("Error marshalling security roles: %s", err.Error()),
			)
			return
		}
		resp.Private.SetKey(ctx, "role", rolesBytes)

		updateUser = *user
	}

	model := convertDataverseFromUserDto(&updateUser, state.DisableDelete.ValueBool())
	state.Id = model.Id
	state.AadId = model.AadId
	state.SecurityRoles = model.SecurityRoles
	state.UserPrincipalName = model.UserPrincipalName
	state.FirstName = model.FirstName
	state.LastName = model.LastName
	state.DisableDelete = model.DisableDelete
	state.BusinessUnitId = model.BusinessUnitId

	tflog.Debug(ctx, fmt.Sprintf("READ: %s_environment with id %s", r.ProviderTypeName, state.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *UserResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hasEnvDataverse, err := r.UserClient.EnvironmentHasDataverse(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Dataverse exist in eviroment %t", hasEnvDataverse))

	addedSecurityRoles, removedSecurityRoles := array.Diff(plan.SecurityRoles, state.SecurityRoles)
	user := userDto{}
	if hasEnvDataverse {
		if len(addedSecurityRoles) > 0 {
			userDto, err := r.UserClient.AddDataverseSecurityRoles(ctx, plan.EnvironmentId.ValueString(), state.Id.ValueString(), addedSecurityRoles)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Client error when adding security roles %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
				return
			}
			user = *userDto
		}
		if len(removedSecurityRoles) > 0 {
			userDto, err := r.UserClient.RemoveDataverseSecurityRoles(ctx, plan.EnvironmentId.ValueString(), state.Id.ValueString(), removedSecurityRoles)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Client error when removing security roles %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
				return
			}
			user = *userDto
		}
	} else {
		err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		}
		if len(addedSecurityRoles) > 0 {
			userDto, err := r.UserClient.AddEnvironmentUserSecurityRoles(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString(), addedSecurityRoles)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Client error when adding security roles %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
				return
			}
			user = *userDto
		}
		if len(removedSecurityRoles) > 0 {
			savedRoles := []securityRoleDto{}
			rolesObj, diag := resp.Private.GetKey(ctx, "role")
			if diag.HasError() {
				resp.Diagnostics.AddError(fmt.Sprintf("Error when updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
				return
			}

			err := json.Unmarshal(rolesObj, &savedRoles)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Error when updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
				return
			}

			userDto, err := r.UserClient.RemoveEnvironmentUserSecurityRoles(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString(), removedSecurityRoles, savedRoles)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Client error when removing security roles %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
				return
			}
			user = *userDto
		}

		rolesBytes, err := json.Marshal(user.SecurityRoles)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error marshalling security roles", fmt.Sprintf("Error marshalling security roles: %s", err.Error()),
			)
			return
		}
		resp.Private.SetKey(ctx, "role", rolesBytes)
	}

	model := convertDataverseFromUserDto(&user, plan.DisableDelete.ValueBool())
	plan.Id = model.Id
	plan.AadId = model.AadId
	req.Plan.SetAttribute(ctx, path.Root("security_roles"), model.SecurityRoles)
	plan.UserPrincipalName = model.UserPrincipalName
	plan.FirstName = model.FirstName
	plan.LastName = model.LastName
	plan.DisableDelete = model.DisableDelete
	plan.BusinessUnitId = model.BusinessUnitId

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	hasEnvDataverse, err := r.UserClient.EnvironmentHasDataverse(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Dataverse exist in eviroment %t", hasEnvDataverse))

	if hasEnvDataverse {
		if state.DisableDelete.ValueBool() {
			err := r.UserClient.DeleteDataverseUser(ctx, state.EnvironmentId.ValueString(), state.Id.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
				return
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Disable delete is set to false. Skipping delete of systemuser with id %s", state.Id.ValueString()))
		}
	} else {

		tflog.Debug(ctx, "11111")

		savedRoles := []securityRoleDto{}
		rolesObj, diag := resp.Private.GetKey(ctx, "role")
		if diag.HasError() {
			return
		}

		tflog.Debug(ctx, "22222")
		err := json.Unmarshal(rolesObj, &savedRoles)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}

		tflog.Debug(ctx, "3333")
		_, err = r.UserClient.RemoveEnvironmentUserSecurityRoles(ctx, state.EnvironmentId.ValueString(), state.AadId.ValueString(), state.SecurityRoles, savedRoles)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
			return
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func validateEnvironmentSecurityRoles(roles []string) error {
	except := array.Except(roles, []string{ROLE_ENVIRONMENT_ADMIN, ROLE_ENVIRONMENT_MAKER})
	if len(except) > 0 {
		return fmt.Errorf("invalid security roles. only '%s' and '%s' are allowed", ROLE_ENVIRONMENT_ADMIN, ROLE_ENVIRONMENT_MAKER)
	}
	return nil
}
