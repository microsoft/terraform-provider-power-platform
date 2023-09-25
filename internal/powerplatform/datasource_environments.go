package powerplatform

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

var (
	_ datasource.DataSource              = &EnvironmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentsDataSource{}
)

func NewEnvironmentsDataSource() datasource.DataSource {
	return &EnvironmentsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environments",
	}
}

type EnvironmentsDataSource struct {
	BapiApiClient    bapi.BapiClientInterface
	ProviderTypeName string
	TypeName         string
}

type EnvironmentsListDataSourceModel struct {
	Environments []EnvironmentDataSourceModel `tfsdk:"environments"`
	Id           types.Int64                  `tfsdk:"id"`
}

type EnvironmentDataSourceModel struct {
	EnvironmentName types.String `tfsdk:"environment_name"`
	DisplayName     types.String `tfsdk:"display_name"`
	Url             types.String `tfsdk:"url"`
	Domain          types.String `tfsdk:"domain"`
	Location        types.String `tfsdk:"location"`
	EnvironmentType types.String `tfsdk:"environment_type"`
	OrganizationId  types.String `tfsdk:"organization_id"`
	SecurityGroupId types.String `tfsdk:"security_group_id"`
	LanguageName    types.Int64  `tfsdk:"language_code"`
	Version         types.String `tfsdk:"version"`
}

func ConvertFromEnvironmentDto(environmentDto models.EnvironmentDto) EnvironmentDataSourceModel {
	return EnvironmentDataSourceModel{
		EnvironmentName: types.StringValue(environmentDto.Name),
		DisplayName:     types.StringValue(environmentDto.Properties.DisplayName),
		Location:        types.StringValue(environmentDto.Location),
		EnvironmentType: types.StringValue(environmentDto.Properties.EnvironmentSku),
		OrganizationId:  types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.ResourceId),
		SecurityGroupId: types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.SecurityGroupId),
		LanguageName:    types.Int64Value(int64(environmentDto.Properties.LinkedEnvironmentMetadata.BaseLanguage)),
		Url:             types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.InstanceURL),
		Domain:          types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.DomainName),
		Version:         types.StringValue(environmentDto.Properties.LinkedEnvironmentMetadata.Version),
	}
}

func (d *EnvironmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *EnvironmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of environments in a tenant",
		MarkdownDescription: "Fetches the list of environments in a tenant",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Placeholder identifier attribute",
				Computed:    true,
			},
			"environments": schema.ListNestedAttribute{
				Description:         "List of environments",
				MarkdownDescription: "List of environments",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"environment_name": schema.StringAttribute{
							MarkdownDescription: "Unique environment name (guid)",
							Description:         "Unique environment name (guid)",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name",
							Description:         "Display name",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							Description:         "Url of the environment",
							MarkdownDescription: "Url of the environment",
							Computed:            true,
						},
						"domain": schema.StringAttribute{
							Description:         "Domain name of the environment",
							MarkdownDescription: "Domain name of the environment",
							Computed:            true,
						},
						"location": schema.StringAttribute{
							Description:         "Location of the environment (europe, unitedstates etc.)",
							MarkdownDescription: "Location of the environment (europe, unitedstates etc.)",
							Computed:            true,
						},
						"environment_type": schema.StringAttribute{
							Description:         "Type of the environment (Sandbox, Production etc.)",
							MarkdownDescription: "Type of the environment (Sandbox, Production etc.)",
							Computed:            true,
						},
						"organization_id": schema.StringAttribute{
							Description:         "Unique organization id (guid)",
							MarkdownDescription: "Unique organization id (guid)",
							Computed:            true,
						},
						"security_group_id": schema.StringAttribute{
							Description:         "Unique security group id (guid)",
							MarkdownDescription: "Unique security group id (guid)",
							Computed:            true,
						},
						"language_code": schema.Int64Attribute{
							Description:         "Unique language LCID (integer)",
							MarkdownDescription: "Unique language LCID (integer)",
							Computed:            true,
						},
						"version": schema.StringAttribute{
							Description:         "Version of the environment",
							MarkdownDescription: "Version of the environment",
							Computed:            true,
						},
						//Not available in BAPI as for now
						// "currency_name": &schema.StringAttribute{
						// 	Description:         "Unique currency name (EUR, USE, GBP etc.)",
						// 	MarkdownDescription: "Unique currency name (EUR, USE, GBP etc.)",
						// 	Computed:            true,
						// },
					},
				},
			},
		},
	}
}

func (d *EnvironmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*PowerPlatformProvider).BapiApi.Client.(bapi.BapiClientInterface)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.BapiApiClient = client
}

func (d *EnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvironmentsListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENTS START: %s", d.ProviderTypeName))

	envs, err := d.BapiApiClient.GetEnvironments(ctx)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, env := range envs {
		e := ConvertFromEnvironmentDto(env)
		state.Environments = append(state.Environments, e)
	}
	state.Id = types.Int64Value(time.Now().Unix())

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENTS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
