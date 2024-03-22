// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var _ resource.Resource = &EnvironmentSettingsResource{}
var _ resource.ResourceWithImportState = &EnvironmentSettingsResource{}

func NewEnvironmentSettingsResource() resource.Resource {
	return &EnvironmentSettingsResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment_settings",
	}
}

type EnvironmentSettingsResource struct {
	EnvironmentSettingClient EnvironmentSettingsClient
	ProviderTypeName         string
	TypeName                 string
}

func (r *EnvironmentSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *EnvironmentSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Power Platform Settings for a given environment.",
		MarkdownDescription: "Manages Power Platform Settings for a given environment. They control various aspects of Power Platform features and behaviors, See [Environment Settings Overview](https://learn.microsoft.com/en-us/power-platform/admin/admin-settings) for more details.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Id",
				Computed:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Environment Id",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_upload_file_size_in_bytes": schema.Int64Attribute{
				Description:         "Maximum file size that can be uploaded to the environment",
				MarkdownDescription: "Maximum file size that can be uploaded to the environment",
				Optional:            true,
			},
			"show_dashboard_cards_in_expanded_state": schema.BoolAttribute{
				Description:         "Show dashboard cards in expanded state",
				MarkdownDescription: "Show dashboard cards in expanded state",
				Optional:            true,
			},
			"plugin_trace_log_setting": schema.StringAttribute{
				Description:         "Plugin trace log setting. Available options: Off, Exception, All",
				MarkdownDescription: "Plugin trace log setting. Available options: Off, Exception, All",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Off", "Exception", "All"),
				},
			},
			"is_audit_enabled": schema.BoolAttribute{
				Description:         "Is audit enabled",
				MarkdownDescription: "Is audit enabled",
				Optional:            true,
			},
			"is_user_access_audit_enabled": schema.BoolAttribute{
				Description:         "Is user access audit enabled",
				MarkdownDescription: "Is user access audit enabled",
				Optional:            true,
			},
			"is_read_audit_enabled": schema.BoolAttribute{
				Description:         "Is read audit enabled",
				MarkdownDescription: "Is read audit enabled",
				Optional:            true,
			},
		},
	}
}

func (r *EnvironmentSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := req.ProviderData.(*api.ProviderClient).Api

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.EnvironmentSettingClient = NewEnvironmentSettingsClient(client)
}

func (r *EnvironmentSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvironmenttSettingsSourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE ENVIRONMENT SETTINGS RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	settingsToUpdate := ConvertFromEnvironmentSettingsModel(plan)

	envSettings, err := r.EnvironmentSettingClient.UpdateEnvironmentSettings(ctx, plan.EnvironmentId.ValueString(), settingsToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment settings", fmt.Sprintf("Error creating environment settings: %s", err.Error()),
		)
		return
	}

	var state = ConvertFromEnvironmentSettingsDto(envSettings)
	state.EnvironmentId = plan.EnvironmentId

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", state.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("CREATE ENVIRONMENT SETTINGS RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvironmenttSettingsSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ ENVIRONMENT SETTINGS RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.EnvironmentId.ValueString() == "" {
		resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
		return
	}

	envSettings, err := r.EnvironmentSettingClient.GetEnvironmentSettings(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.ProviderTypeName), err.Error())
		return
	}

	var newState = ConvertFromEnvironmentSettingsDto(envSettings)
	newState.EnvironmentId = state.EnvironmentId

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	tflog.Debug(ctx, fmt.Sprintf("READ ENVIRONMENT SETTINGS RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvironmenttSettingsSourceModel
	var state EnvironmenttSettingsSourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE ENVIRONMENT SETTINGS RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.EnvironmentId.ValueString() == "" {
		resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
		return
	}

	envSettingsToUpdate := ConvertFromEnvironmentSettingsModel(plan)

	environmentSettings, err := r.EnvironmentSettingClient.UpdateEnvironmentSettings(ctx, plan.EnvironmentId.ValueString(), envSettingsToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating environment settings", fmt.Sprintf("Error updating environment settings: %s", err.Error()),
		)
		return
	}

	plan = ConvertFromEnvironmentSettingsDto(environmentSettings)
	plan.EnvironmentId = state.EnvironmentId

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("UPDATE ENVIRONMENT SETTINGS RESOURCE END: %s", r.ProviderTypeName))
}

func (r *EnvironmentSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *EnvironmentSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
