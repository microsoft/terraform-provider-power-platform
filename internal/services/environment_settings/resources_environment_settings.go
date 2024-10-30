// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var _ resource.Resource = &EnvironmentSettingsResource{}
var _ resource.ResourceWithImportState = &EnvironmentSettingsResource{}

func NewEnvironmentSettingsResource() resource.Resource {
	return &EnvironmentSettingsResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_settings",
		},
	}
}

func (r *EnvironmentSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	r.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = r.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (r *EnvironmentSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Manages Power Platform Settings for a given environment.",
		MarkdownDescription: "Manages Power Platform Settings for a given environment. They control various aspects of Power Platform features and behaviors, See [Environment Settings Overview](https://learn.microsoft.com/power-platform/admin/admin-settings) for more details.  While this resource provides a limited set of settings, many of the settings in an environment are stored as Dataverse records and can be managed using `powerplatform_data_record` resource.  See the [data record resource documentation](./data_record) for examples of how to manage more environment settings like business units, roles, and more.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
			}),
			"id": schema.StringAttribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				Description:         "Environment Id",
				MarkdownDescription: "Environment Id",
				Required:            true,
			},
			"audit_and_logs": schema.SingleNestedAttribute{
				Description:         "Audit and Logs",
				MarkdownDescription: "Audit and Logs",
				Optional:            true, Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"plugin_trace_log_setting": schema.StringAttribute{
						Description:         "Plugin trace log setting. Available options: Off, Exception, All",
						MarkdownDescription: "Plugin trace log setting. Available options: Off, Exception, All. See [Plugin Trace Log Settings Overview](https://learn.microsoft.com/power-apps/developer/data-platform/logging-tracing) for more details.",
						Optional:            true, Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.OneOf("Off", "Exception", "All"),
						},
					},
					"audit_settings": schema.SingleNestedAttribute{
						Description:         "Audit Settings",
						MarkdownDescription: "Audit Settings. See [Audit Settings Overview](https://learn.microsoft.com/power-platform/admin/system-settings-dialog-box-auditing-tab) for more details.",
						Optional:            true, Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"is_audit_enabled": schema.BoolAttribute{
								Description:         "Is audit enabled",
								MarkdownDescription: "Is audit enabled",
								Optional:            true, Computed: true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"is_user_access_audit_enabled": schema.BoolAttribute{
								Description:         "Is user access audit enabled",
								MarkdownDescription: "Is user access audit enabled",
								Optional:            true, Computed: true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"is_read_audit_enabled": schema.BoolAttribute{
								Description:         "Is read audit enabled",
								MarkdownDescription: "Is read audit enabled",
								Optional:            true, Computed: true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"email": schema.SingleNestedAttribute{
				Description:         "Email",
				MarkdownDescription: "Email",
				Optional:            true, Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"email_settings": schema.SingleNestedAttribute{
						Description:         "Email Settings",
						MarkdownDescription: "Email Settings. See [Email Settings Overview](https://learn.microsoft.com/power-platform/admin/settings-email) for more details.",
						Optional:            true, Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"max_upload_file_size_in_bytes": schema.Int64Attribute{
								Description:         "Maximum file size that can be uploaded to the environment",
								MarkdownDescription: "Maximum file size that can be uploaded to the environment",
								Optional:            true, Computed: true,
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"product": schema.SingleNestedAttribute{
				Description:         "Product",
				MarkdownDescription: "Product",
				Optional:            true, Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"behavior_settings": schema.SingleNestedAttribute{
						Description:         "Behavior Settings",
						MarkdownDescription: "Behavior Settings.See [Behavior Settings Overview](https://learn.microsoft.com/power-platform/admin/settings-behavior) for more details.",
						Optional:            true, Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"show_dashboard_cards_in_expanded_state": schema.BoolAttribute{
								Description:         "Show dashboard cards in expanded state",
								MarkdownDescription: "Show dashboard cards in expanded state",
								Optional:            true, Computed: true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"features": schema.SingleNestedAttribute{
						Description:         "Features",
						MarkdownDescription: "Features. See [Features Overview](https://learn.microsoft.com/power-platform/admin/settings-features) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"power_apps_component_framework_for_canvas_apps": schema.BoolAttribute{
								Description:         "Power Apps component framework for canvas apps",
								MarkdownDescription: "Power Apps component framework for canvas apps",
								Optional:            true, Computed: true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *EnvironmentSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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

	r.EnvironmentSettingClient = newEnvironmentSettingsClient(client)
}

func (r *EnvironmentSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settingsToUpdate, err := convertFromEnvironmentSettingsModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting environment settings model", err.Error())
		return
	}

	dvExits, err := r.EnvironmentSettingClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
		return
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
		return
	}

	envSettings, err := r.EnvironmentSettingClient.UpdateEnvironmentSettings(ctx, plan.EnvironmentId.ValueString(), *settingsToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment settings", fmt.Sprintf("Error creating environment settings: %s", err.Error()),
		)
		return
	}

	var state = convertFromEnvironmentSettingsDto[EnvironmentSettingsResourceModel](envSettings, plan.Timeouts)
	state.Id = plan.EnvironmentId
	state.EnvironmentId = plan.EnvironmentId

	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID %s", state.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EnvironmentSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state EnvironmentSettingsResourceModel
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

	var newState = convertFromEnvironmentSettingsDto[EnvironmentSettingsResourceModel](envSettings, state.Timeouts)
	newState.Id = state.EnvironmentId
	newState.EnvironmentId = state.EnvironmentId

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *EnvironmentSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan EnvironmentSettingsResourceModel
	var state EnvironmentSettingsResourceModel
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

	envSettingsToUpdate, err := convertFromEnvironmentSettingsModel(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting environment settings model", err.Error())
		return
	}

	environmentSettings, err := r.EnvironmentSettingClient.UpdateEnvironmentSettings(ctx, plan.EnvironmentId.ValueString(), *envSettingsToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating environment settings", fmt.Sprintf("Error updating environment settings: %s", err.Error()),
		)
		return
	}

	plan = convertFromEnvironmentSettingsDto[EnvironmentSettingsResourceModel](environmentSettings, plan.Timeouts)
	plan.Id = state.EnvironmentId
	plan.EnvironmentId = state.EnvironmentId

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EnvironmentSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	// Do nothing on purpose
}

func (r *EnvironmentSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	resource.ImportStatePassthroughID(ctx, path.Root("environment_id"), req, resp)
}
