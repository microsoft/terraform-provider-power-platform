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
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (d *EnvironmentSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.EnvironmentSettingsClient = newEnvironmentSettingsClient(client.Api)
}

func (d *EnvironmentSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	var state EnvironmentSettingsDataSourceModel

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

	state = convertFromEnvironmentSettingsDto[EnvironmentSettingsDataSourceModel](envSettings, state.Timeouts)

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
		MarkdownDescription: "Power Platform Environment Settings Data Source",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: false,
				Update: false,
				Delete: false,
				Read:   false,
			}),
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Unique environment id (guid)",
				Required:            true,
			},
			"audit_and_logs": schema.SingleNestedAttribute{
				MarkdownDescription: "Audit and Logs",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"plugin_trace_log_setting": schema.StringAttribute{
						MarkdownDescription: "Plugin trace log setting. Available options: Off, Exception, All. See [Plugin Trace Log Settings Overview](https://learn.microsoft.com/power-apps/developer/data-platform/logging-tracing) for more details.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("Off", "Exception", "All"),
						},
					},
					"audit_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Audit Settings. See [Audit Settings Overview](https://learn.microsoft.com/power-platform/admin/system-settings-dialog-box-auditing-tab) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"is_audit_enabled": schema.BoolAttribute{
								MarkdownDescription: "Is audit enabled",
								Optional:            true,
							},
							"is_user_access_audit_enabled": schema.BoolAttribute{
								MarkdownDescription: "Is user access audit enabled",
								Optional:            true,
							},
							"is_read_audit_enabled": schema.BoolAttribute{
								MarkdownDescription: "Is read audit enabled",
								Optional:            true,
							},
							"log_retention_period_in_days": schema.Int32Attribute{
								MarkdownDescription: "Retain these logs for. See [Start/stop auditing for an environment and set retention policy](https://learn.microsoft.com/power-platform/admin/manage-dataverse-auditing#startstop-auditing-for-an-environment-and-set-retention-policy) You can set a retention period for how long audit logs are kept in an environment. Under Retain these logs for, choose the period of time you wish to retain the logs.",
								Optional:            true,
							},
						},
					},
				},
			},
			"email": schema.SingleNestedAttribute{
				MarkdownDescription: "Email",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"email_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Email Settings. See [Email Settings Overview](https://learn.microsoft.com/power-platform/admin/settings-email) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"max_upload_file_size_in_bytes": schema.Int64Attribute{
								MarkdownDescription: "Maximum file size that can be uploaded to the environment",
								Optional:            true,
							},
						},
					},
				},
			},
			"product": schema.SingleNestedAttribute{
				MarkdownDescription: "Product",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"behavior_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Behavior Settings.See [Behavior Settings Overview](https://learn.microsoft.com/power-platform/admin/settings-behavior) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"show_dashboard_cards_in_expanded_state": schema.BoolAttribute{
								MarkdownDescription: "Show dashboard cards in expanded state",
								Optional:            true,
							},
						},
					},
					"features": schema.SingleNestedAttribute{
						MarkdownDescription: "Features. See [Features Overview](https://learn.microsoft.com/power-platform/admin/settings-features) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"power_apps_component_framework_for_canvas_apps": schema.BoolAttribute{
								MarkdownDescription: "Power Apps component framework for canvas apps",
								Optional:            true,
							},
						},
					},
					"security": schema.SingleNestedAttribute{
						MarkdownDescription: "Security. See [Security Overview](https://learn.microsoft.com/en-us/power-platform/admin/settings-privacy-security) for more details.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enable_ip_based_cookie_binding": schema.BoolAttribute{
								MarkdownDescription: "Enable IP based cookie binding",
								Optional:            true,
							},
							"enable_ip_based_firewall_rule": schema.BoolAttribute{
								MarkdownDescription: "Enable IP based firewall rule",
								Optional:            true,
							},
							"allowed_ip_range_for_firewall": schema.SetAttribute{
								MarkdownDescription: "Allowed IP range for firewall",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"allowed_service_tags_for_firewall": schema.SetAttribute{
								MarkdownDescription: "Allowed service tags for firewall",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"allow_application_user_access": schema.BoolAttribute{
								MarkdownDescription: "Allow application user access",
								Optional:            true,
							},
							"allow_microsoft_trusted_service_tags": schema.BoolAttribute{
								MarkdownDescription: "Allow Microsoft trusted service tags",
								Optional:            true,
							},
							"enable_ip_based_firewall_rule_in_audit_mode": schema.BoolAttribute{
								MarkdownDescription: "Enable IP based firewall rule in audit mode",
								Optional:            true,
							},
							"reverse_proxy_ip_addresses": schema.SetAttribute{
								MarkdownDescription: "Reverse proxy IP addresses",
								ElementType:         types.StringType,
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
