// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &EnvironmentSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentSettingsDataSource{}
)

func NewEnvironmentSettingsDataSource() *EnvironmentSettingsDataSource {
	return &EnvironmentSettingsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_settings",
		},
	}
}

type EnvironmentSettingsDataSource struct {
	helpers.TypeInfo
	EnvironmentSettingsClient client
}

func (d *EnvironmentSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
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
	d.EnvironmentSettingsClient = newEnvironmentSettingsClient(client)
}

func (d *EnvironmentSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	var state EnvironmentSettingsSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT SETTINGS START: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.EnvironmentId.ValueString() == "" {
		resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
		return
	}

	dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
		return
	}

	envSettings, err := d.EnvironmentSettingsClient.GetEnvironmentSettings(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	state = convertFromEnvironmentSettingsDto(envSettings, state.Timeouts)
	state.Id = state.EnvironmentId

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT SETTINGS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *EnvironmentSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Power Platform Environment Settings Data Source",
		MarkdownDescription: "Power Platform Environment Settings Data Source. Power Platform Settings are configuration options that apply to a specific environment. They control various aspects of Power Platform features and behaviors, See [Environment Settings Overview](https://learn.microsoft.com/power-platform/admin/admin-settings) for more details.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: false,
				Update: false,
				Delete: false,
				Read:   false,
			}),
			"id": schema.StringAttribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				Description:         "Unique environment id (guid)",
				MarkdownDescription: "Unique environment id (guid)",
				Required:            true,
			},
			"audit_and_logs": schema.SingleNestedAttribute{
				Description:         "Audit and Logs",
				MarkdownDescription: "Audit and Logs",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"plugin_trace_log_setting": schema.StringAttribute{
						Description:         "Plugin trace log setting. Available options: Off, Exception, All",
						MarkdownDescription: "Plugin trace log setting. Available options: Off, Exception, All. See [Plugin Trace Log Settings Overview](https://learn.microsoft.com/power-apps/developer/data-platform/logging-tracing) for more details.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("Off", "Exception", "All"),
						},
					},
					"audit_settings": schema.SingleNestedAttribute{
						Description:         "Audit Settings",
						MarkdownDescription: "Audit Settings. See [Audit Settings Overview](https://learn.microsoft.com/power-platform/admin/system-settings-dialog-box-auditing-tab) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
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
					},
				},
			},
			"email": schema.SingleNestedAttribute{
				Description:         "Email",
				MarkdownDescription: "Email",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"email_settings": schema.SingleNestedAttribute{
						Description:         "Email Settings",
						MarkdownDescription: "Email Settings. See [Email Settings Overview](https://learn.microsoft.com/power-platform/admin/settings-email) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"max_upload_file_size_in_bytes": schema.Int64Attribute{
								Description:         "Maximum file size that can be uploaded to the environment",
								MarkdownDescription: "Maximum file size that can be uploaded to the environment",
								Optional:            true,
							},
						},
					},
				},
			},
			"product": schema.SingleNestedAttribute{
				Description:         "Product",
				MarkdownDescription: "Product",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"behavior_settings": schema.SingleNestedAttribute{
						Description:         "Behavior Settings",
						MarkdownDescription: "Behavior Settings.See [Behavior Settings Overview](https://learn.microsoft.com/power-platform/admin/settings-behavior) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"show_dashboard_cards_in_expanded_state": schema.BoolAttribute{
								Description:         "Show dashboard cards in expanded state",
								MarkdownDescription: "Show dashboard cards in expanded state",
								Optional:            true,
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
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}

func (d *EnvironmentSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
