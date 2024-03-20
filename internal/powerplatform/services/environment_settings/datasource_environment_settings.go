// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &EnvironmentSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentSettingsDataSource{}
)

func NewEnvironmentSettingsDataSource() *EnvironmentSettingsDataSource {
	return &EnvironmentSettingsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment_settings",
	}
}

type EnvironmentSettingsDataSource struct {
	EnvironmentSettingsClient EnvironmentSettingsClient
	ProviderTypeName          string
	TypeName                  string
}

type EnvironmentSettingsDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	EnvironmentId                     types.String `tfsdk:"environment_id"`
	MaxUploadFileSize                 types.Int64  `tfsdk:"max_upload_file_size"`
	ShowDashboardCardsInExpandedState types.Bool   `tfsdk:"show_dashboard_cards_in_expanded_state"`
	PluginTraceLogSetting             types.String `tfsdk:"plugin_trace_log_setting"`
	IsAuditEnabled                    types.Bool   `tfsdk:"is_audit_enabled"`
	IsUserAccessAuditEnabled          types.Bool   `tfsdk:"is_user_access_audit_enabled"`
	IsReadAuditEnabled                types.Bool   `tfsdk:"is_read_audit_enabled"`
}

func (d *EnvironmentSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.EnvironmentSettingsClient = NewEnvironmentSettingsClient(client)

}

func (d *EnvironmentSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvironmentSettingsDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT SETTINGS START: %s", d.ProviderTypeName))

	if state.EnvironmentId.ValueString() == "" {
		resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
		return
	}

	envSettings, err := d.EnvironmentSettingsClient.GetEnvironmentSettings(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	state.Id = types.StringValue(uuid.New().String())
	state.MaxUploadFileSize = types.Int64Value(envSettings.MaxUploadFileSize)
	state.ShowDashboardCardsInExpandedState = types.BoolValue(envSettings.BoundDashboardDefaultCardExpanded)

	state.PluginTraceLogSetting = types.StringValue(d.getPluginTraceLogSetting(envSettings.PluginTraceLogSetting))
	state.IsAuditEnabled = types.BoolValue(envSettings.IsAuditEnabled)
	state.IsUserAccessAuditEnabled = types.BoolValue(envSettings.IsUserAccessAuditEnabled)
	state.IsReadAuditEnabled = types.BoolValue(envSettings.IsReadAuditEnabled)

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT SETTINGS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *EnvironmentSettingsDataSource) getPluginTraceLogSetting(value int64) string {
	switch value {
	case 0:
		return "Off"
	case 1:
		return "Exception"
	case 2:
		return "All"
	default:
		panic(fmt.Sprintf("Invalid value for pluginTraceLogSetting: %d", value))
	}
}

func (d *EnvironmentSettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Power Platform Tenant Settings Data Source",
		MarkdownDescription: "Power Platform Tenant Settings Data Source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Id",
				Computed:    true,
			},
			"environment_id": schema.StringAttribute{
				Description:         "Unique environment id (guid)",
				MarkdownDescription: "Unique environment id (guid)",
				Required:            true,
			},
			"max_upload_file_size": schema.Int64Attribute{
				Description:         "Maximum file size that can be uploaded to the environment",
				MarkdownDescription: "Maximum file size that can be uploaded to the environment",
				Computed:            true,
			},
			"show_dashboard_cards_in_expanded_state": schema.BoolAttribute{
				Description:         "Show dashboard cards in expanded state",
				MarkdownDescription: "Show dashboard cards in expanded state",
				Computed:            true,
			},
			"plugin_trace_log_setting": schema.StringAttribute{
				Description:         "Plugin trace log setting. Available options: Off, Exception, All",
				MarkdownDescription: "Plugin trace log setting. Available options: Off, Exception, All",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Off", "Exception", "All"),
				},
			},
			"is_audit_enabled": schema.BoolAttribute{
				Description:         "Is audit enabled",
				MarkdownDescription: "Is audit enabled",
				Computed:            true,
			},
			"is_user_access_audit_enabled": schema.BoolAttribute{
				Description:         "Is user access audit enabled",
				MarkdownDescription: "Is user access audit enabled",
				Computed:            true,
			},
			"is_read_audit_enabled": schema.BoolAttribute{
				Description:         "Is read audit enabled",
				MarkdownDescription: "Is read audit enabled",
				Computed:            true,
			},
		},
	}
}

func (d *EnvironmentSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}
